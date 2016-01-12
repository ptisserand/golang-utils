#include <stdio.h>
#include "foo.h"

int main(int argc, char** argv) {
    fprintf(stderr, "Foo: %d\n", foo);
    foo = 42;
    fprintf(stderr, "Foo: %d\n", foo);
    fprintf(stderr, "Billy: [");
    for (int i = 0; i < 4; ++i) {
	fprintf(stderr, " 0x%.2x", billy[i]); 
    }
    fprintf(stderr, " ]\n");
    fprintf(stderr, "Bob: [");
    for (int i = 0; i < 4; ++i) {
	fprintf(stderr, " 0x%.2x", bob[i]); 
    }
    fprintf(stderr, " ]\n");

    uint8_t* pp = &(yop);
    fprintf(stderr, "Yop: [");
    for (int i = 0; i < 8; ++i) {
	fprintf(stderr, " 0x%.2x", pp[i]); 
    }
    fprintf(stderr, " ]\n");
    
    return 0;
}
