#include "gradient.hpp"
#include "random.hpp"
#include <cstdlib>

int *gradient::array() {
  int *colors[(this->capacity() * 3)];
  for (int i = 0; i < this->capacity(); i++) {
    colors[i] = &this->grad[i].R;
    colors[i + 1] = &this->grad[i].G;
    colors[i + 2] = &this->grad[i].B;
  }
  return *colors;
}

int *gradient::arrayWithHeight(int HEIGHT) {
  int *colors[(this->capacity() * 3) * HEIGHT];
  for (int n = 0; n < HEIGHT; n++) {
    for (int i = 0; i < this->capacity(); i++) {
      colors[n * i] = &this->grad[i].R;
      colors[n * i + 1] = &this->grad[i].G;
      colors[n * i + 2] = &this->grad[i].B;
    }
  }
  return *colors;
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