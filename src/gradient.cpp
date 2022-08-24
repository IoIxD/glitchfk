#include "gradient.hpp"
#include "random.hpp"
#include <cstdlib>
#include <cstdio>

using std::printf;

int* gradient::array(int width, int height) {
  int *colors = (int*)malloc(width*height*12);
  for (int y = 0; y < height; y++) {
    //printf("%d\n",y);
    for (int x = 0; x < width; x+=3) {
      //printf("%d\n",x);
      int point = (y*x)+x;
      colors[point] = this->grad[x].R;
      colors[point + 1] = this->grad[x].G;
      colors[point + 2] = this->grad[x].B;
    }
  }
  return colors;
}

gradient LinearGradient(color color1, color color2, int width) {
  gradient grad;
  for (int x = 0; x < width; x++) {
    color newColor{
        (x * color1.R + x * (color2.R - color1.R)) / 1024,
        (x * color1.G + x * (color2.G - color1.G)) / 1024,
        (x * color1.B + x * (color2.B - color1.B)) / 1024,
    };
    grad.append(newColor);
  };
  return grad;
};

gradient RandomLinearGradient(int width) {
  srand(nanosecondSeed());
  color Color1{
      (int)(random() % 255),
      (int)(random() % 255),
      (int)(random() % 255),
  };
  srand(nanosecondSeed());
  color Color2{
      (int)(random() % 255),
      (int)(random() % 255),
      (int)(random() % 255),
  };
  return LinearGradient(Color1, Color2, width);
}