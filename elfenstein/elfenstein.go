package main

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"gopkg.in/ini.v1"
)

type RElfData struct {
	Name string
	Offset uint64 // Offset in ELF
	Size uint64
}

func display_relf_data(adata RElfData) {
	fmt.Printf("Name: %s\n", adata.Name)
	fmt.Printf("Offset: 0x%.8x\n", adata.Offset)
	fmt.Printf("Size: 0x%.8x\n", adata.Size)
}

func replace_elf_data(adata RElfData, outf *os.File, fname string) (error) {
	fstat, err := os.Lstat(fname)
	if err != nil {
		return err
	}
	var size uint64 = uint64(fstat.Size())
	if size > adata.Size {
		return fmt.Errorf("%s size (%d) is greater than corresponding ELF data (%d)\n", fname, size, adata.Size)
	}
	data := make([]byte, size)
	af, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer af.Close()
	n, err := af.Read(data)
	if err != nil {
		return err
	}
	n, err = outf.WriteAt(data[:n], int64(adata.Offset))
	if err != nil {
		return err
	}
	return nil
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func find_elf_section(aelf *elf.File, name string) (*RElfData, error) {
	var found bool = false
	ret := new(RElfData)
	for _, asect := range aelf.Sections {
		if asect.SectionHeader.Name == name {
			ret.Name = name
			ret.Offset = asect.SectionHeader.Offset
			ret.Size = asect.SectionHeader.Size
			found = true
			break;
		}
	}
	if found {
		return ret, nil
	} else {
		return nil, fmt.Errorf("Section '%s' not present", name)
	}
}

func find_elf_symbol(aelf *elf.File, name string) (*RElfData, error) {
	var found bool = false
	ret := new(RElfData)
	symbols, _ := aelf.Symbols()
	for _, asym := range symbols {
		if asym.Name == name {
			ret.Name = name
			ret.Size = asym.Size
			asect := aelf.Sections[asym.Section]
			ret.Offset = asect.SectionHeader.Offset + (asym.Value - asect.SectionHeader.Addr)
			found = true
			break
		}
	}
	if found {
		return ret, nil
	} else {
		return nil, fmt.Errorf("Symbol '%s' not found", name)
	}
}

func display_elf_section_header(aelf *elf.File, name string) {
	asect, err := find_elf_section(aelf, name)
	if asect != nil {
		display_relf_data(*asect)
	} else {
		panic(err)
	}
}

func display_elf_symbol(aelf *elf.File, name string) {
	asym, err := find_elf_symbol(aelf, name)
	if asym != nil {
		display_relf_data(*asym)
	} else {
		panic(err)
	}
}

func copy_file(src string, dst string) (error){
	fsrc, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fdst.Close()

	size, err := io.Copy(fdst, fsrc)
	if err != nil {
		return err
	}
	srcStat, err := os.Lstat(src)
	if size != srcStat.Size() {
		return fmt.Errorf("%s %d/%d copied", src, size, srcStat.Size())
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: elfenstein config_file elf_file [out_elf_file]")
		os.Exit(1)
	}

	cfg, _ := ini.Load(os.Args[1])


	var outf *os.File
	if len(os.Args) == 4 {
		copy_file(os.Args[2], os.Args[3])
		outf, _ = os.OpenFile(os.Args[3], os.O_WRONLY, 0666)
		defer outf.Close()
	} else {
		outf = nil
	}
	
	_elf, err := elf.Open(os.Args[2])
	check(err)

	fmt.Println("####### CONFIG FILE ITEM in 'sections' ");	
	asect := cfg.Section("sections")
	for _, ii := range asect.Keys() {
		display_elf_section_header(_elf, ii.Name())
		if ii.Value() != "" {
			fmt.Printf("Value: %s\n", ii.Value())
			if outf != nil {
				edata, err := find_elf_section(_elf, ii.Name())
				if edata != nil {
					err = replace_elf_data(*edata, outf, ii.Value())
					check(err)
				}
			}			
		}
	}

	fmt.Println("####### CONFIG FILE ITEM in 'symbols' ");	
	asect = cfg.Section("symbols")
	for _, ii := range asect.Keys() {
		display_elf_symbol(_elf, ii.Name())
		if ii.Value() != "" {
			fmt.Printf("Value: %s\n", ii.Value())
			if outf != nil {
				edata, err := find_elf_symbol(_elf, ii.Name())
				if edata != nil {
					err = replace_elf_data(*edata, outf, ii.Value())
					check(err)
				}
			}			
		}		
	}
	_elf.Close()
}
