#include <chrono>
#include <vector>

using std::vector;
using namespace std::chrono;

typedef class color {
public:
  int R;
  int G;
  int B;
} color;

typedef class gradient {
private:
  vector<color> grad;

public:
  class gradient operator^(gradient grad2) const {
    gradient newGrad;
    for (int x = 0; x < grad.capacity(); x++) {
      color newColor{
          grad[x].R ^ grad2[x].R,
          grad[x].G ^ grad2[x].G,
          grad[x].B ^ grad2[x].B,
      };
      newGrad.grad.push_back(newColor);
    };
    return newGrad;
  }

  int* array(int width, int height);

  void append(color col) { grad.push_back(col); }
  int capacity() { return grad.capacity(); }

  color operator[](int i) const { return grad[i]; };
} gradient;

gradient LinearGradient(color color1, color color2, int width);
gradient RandomLinearGradient(int width);