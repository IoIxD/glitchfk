#include <cstddef>
#include <cstdio>
#include <fcntl.h>
#include <iostream>
#include <ranges>
#include <string>
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
  size_t size = (WIDTH*HEIGHT)*12;
  auto image = VImage::new_from_memory_steal(finalgrad.array(WIDTH,HEIGHT), size, WIDTH, HEIGHT,
                                       3, VIPS_FORMAT_INT);

  image.write_to_file("test.png");

  vips_shutdown();

  return 0;
}
