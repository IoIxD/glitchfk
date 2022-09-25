#include <chrono>
#include <vector>

using std::vector;
using namespace std::chrono;

typedef class color {
public:
  double R;
  double G;
  double B;
} color;

typedef class gradient {
private:
  vector<color> grad;

public:
  unsigned int type;
  unsigned int width;
  unsigned int height;

  void append(color col) { grad.push_back(col); }
  int capacity() { return grad.capacity(); }

  int *operator^(gradient grad2);
  color operator[](int i) const;

  int *array();

  gradient();

} gradient;

gradient LinearGradient(color color1, color color2, double width);
gradient RandomLinearGradient(int width, int height);
gradient RandomHorizontalGradient(int width, int height);
gradient RandomVerticalGradient(int width, int height);