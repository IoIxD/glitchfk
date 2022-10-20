
// The structs and functions in this file allow us, in theory,
// to have an array and a vector that can be interchangably used,
// as well as the ability to turn a vector into an array if it is of a common
// size.

// CURRENTLY THIS IS UNUSED but i'm keeping it here because if i can find out how to
// hack a vec into an array in a timely manner then a big speed boost can be had.

use crate::tinier_gradient::rgb::RGB;

#[derive(Clone)]
#[derive(Debug)]
pub enum ArrOrVec {
    Vec(Vec<RGB>),
    Arr(AcceptableVecSize),
}

impl ArrOrVec {
    pub fn len(&self) -> usize {
        match self {
            ArrOrVec::Vec(a) => a.len(),
            ArrOrVec::Arr(a) => {
                match a {
                    AcceptableVecSize::Type1(..) => 1,
                    AcceptableVecSize::Type640480(..) => 640*480,
                    AcceptableVecSize::Type800600(..) => 800*600,
                }
            }
        }
    }
    pub fn at(&self,i: usize) -> &RGB {
        match self {
            ArrOrVec::Vec(a) => &a.get(i).unwrap(),
            ArrOrVec::Arr(a) => {
                match a {
                    AcceptableVecSize::Type1(a) => &a.get(i).unwrap(),
                    AcceptableVecSize::Type640480(a) => &a.get(i).unwrap(),
                    AcceptableVecSize::Type800600(a) => &a.get(i).unwrap(),
                }
            }
        }
    }
}

impl std::ops::Index<usize> for ArrOrVec {
    type Output = RGB;

    fn index(&self, i: usize) -> &Self::Output {
        self.at(i)
    }
}

#[derive(Clone)]
#[derive(Debug)]
pub enum AcceptableVecSize {
    Type1([RGB; 8]),
    Type640480([RGB; 640*480]),
    Type800600([RGB; 800*600]),
}

pub mod macros {
    #[allow(unused_macros)]

    // macro for trying to stuff a vector with a known size to an array.
    macro_rules! vec_to_array {
        ($elem:expr; $type:ty; $default:expr; $w:expr; $h:expr) => {
            {
                let mut arr: [$type; $w*$h] = [$default; $w*$h];
                for n in 0..$w*$h {
                    arr[n] = *$elem.at(n);
                }
                arr
            }
        };
    }
    //pub(crate) use vec_to_array;
}