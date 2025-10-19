import os

def synthesize_header(path: str, name: str) -> str:
    bytes = ""
    with open(path, 'rb') as f:
        while (byte := f.read(1)):
            bytes += f"0x{byte.hex()}, "
    return f"""#ifndef {name.upper()}_H
#define {name.upper()}_H
#include <stdint.h>
uint8_t {name}Bin[] = {{ {bytes} }};
#endif"""

cobalt_header = synthesize_header("./assets/Cobalt.dll", "cobalt")
with open("src/cobalt.h", 'w') as f:
    f.write(cobalt_header)