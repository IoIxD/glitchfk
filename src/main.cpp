#include <cstddef>
#include <cstdio>
#include <fcntl.h>
#include <iostream>
#include <ranges>
#include <string>
#include <vips/vips.h>
#include <vips/vips8>

#define VIPS_DEBUG
#define VIPS_DEBUG_VERBOSE

#include "gradient.hpp"

using std::printf;

using namespace vips;

using vips::VImage;

const int WIDTH = 1024;
const int HEIGHT = 768;

int main(int argc, const char *argv[]) {
  if (vips_init(argv[0]) != 0) {
    vips_error_exit(nullptr);
  }

  // generate two gradients
  gradient grad1 = RandomLinearGradient(WIDTH);
  gradient grad2 = RandomLinearGradient(WIDTH);

  // generate a xor'd gradient
  gradient finalgrad = grad1 ^ grad2;

  // write that xor'd gradient to an image file
  size_t size;
  auto image = vips_image_new_from_memory(finalgrad.array(), size, WIDTH, 1, 3,
                                          VIPS_FORMAT_INT);

  int err = vips_image_write_to_file(image, "test.png", NULL);

  if (err != 0) {
    printf("%s", vips_error_buffer());
  }

  vips_shutdown();

  return 0;
}
