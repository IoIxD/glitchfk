#include "gradient.hpp"
#include "random.hpp"
#include <cstdio>
#include <cstdlib>
#include <stdexcept>

int *gradient::array() {
  int type = this->type;
  int width = this->width;
  int height = this->height;
  int *colors = (int *)malloc(width * height * 12);
  for (int y = 0; y < height; y++) {
    for (int x = 0; x < width; x++) {
      int point;
      switch (type) {
      case 0:
        point = y * width * 3 + x * 3;
        break;
      case 1:
        point = x * height * 3 + y * 3;
        break;
      }
      colors[point] = this->grad[x].R;
      colors[point + 1] = this->grad[x].G;
      colors[point + 2] = this->grad[x].B;
    }
  }
  return colors;
}

gradient LinearGradient(color color1, color color2, double length) {
  gradient grad;
  for (double x = 0; x < length; x++) {
    color newColor{
        (color2.R - color1.R) * (length / x) + color1.R,
        (color2.G - color1.G) * (length / x) + color1.G,
        (color2.B - color1.B) * (length / x) + color1.B,
    };
    grad.append(newColor);
  };
  return grad;
};

gradient RandomLinearGradient(int width, int height) {
  srand(nanosecondSeed());
  color Color1{
      (double)(random() % 255),
      (double)(random() % 255),
      (double)(random() % 255),
  };
  srand(nanosecondSeed());
  color Color2{
      (double)(random() % 255),
      (double)(random() % 255),
      (double)(random() % 255),
  };
  int length;
  if (height < width) {
    length = width;
  } else {
    length = height;
  }
  gradient grad = LinearGradient(Color1, Color2, length);
  grad.width = width;
  grad.height = height;
  return grad;
}

gradient RandomHorizontalGradient(int width, int height) {
  gradient grad = RandomLinearGradient(width, height);
  grad.type = 0;
  return grad;
}

gradient RandomVerticalGradient(int width, int height) {
  gradient grad = RandomLinearGradient(width, height);
  grad.type = 1;
  return grad;
}

// operators

int *gradient::operator^(gradient grad2) {
  int *grad1arr = this->array();
  int *grad2arr = this->array();

  gradient newGrad;

  if (sizeof(grad1arr) != sizeof(grad2arr)) {
    throw std::invalid_argument("Gradient sizes don't match.");
  }

  int width1 = this->width;
  int height1 = this->height;
  int width2 = grad2.width;
  int height2 = grad2.height;
  int type1 = this->type;
  int type2 = grad2.type;

  newGrad.width = width1;
  newGrad.height = height1;

  for (int y = 0; y < height; y++) {
    for (int x = 0; x < width; x++) {
      int p1, p2;
      switch (type) {
      case 0:
        p1 = y * width1 * 3 + x * 3;
        p2 = y * width2 * 3 + x * 3;
        break;
      case 1:
        p1 = x * height1 * 3 + y * 3;
        p2 = x * height2 * 3 + y * 3;
        break;
      }
      color newColor{
          (double)((int)this->grad[p1].R ^ (int)grad2[p2].R),
          (double)((int)this->grad[p1].G ^ (int)grad2[p2].G),
          (double)((int)this->grad[p1].B ^ (int)grad2[p2].B),
      };
      newGrad.grad.push_back(newColor);
    };
  };

  return newGrad.array();
}

color gradient::operator[](int i) const { return grad[i]; };

gradient::gradient() {
  this->type = 0;
  this->width = 1;
  this->height = 1;
}
