#include <stdint.h>
#include "foo.h"

uint32_t foo __attribute__((section(".mydata"))) = 10;
uint8_t billy[4] __attribute__((section(".mydata"))) = {
    0x11, 0x22, 0x33, 0x44
};
uint64_t yop __attribute__((section(".my2nddata"))) = 0xCACADEADBEEFCACA;
uint8_t bob[4] __attribute__((section(".my2nddata"))) = {
    0xAA, 0xBB, 0xCC, 0xDD
};

