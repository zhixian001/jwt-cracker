#pragma once

#include "mpfr.h"
#include "gmp.h"

// MPFR Guide: https://www.mpfr.org/mpfr-current/mpfr.html#index-mpfr_005frnd_005ft
 
/**
 * @brief Calculate the natural logarithm of an (big) integer
 * 
 * @param input_int arbitrarily (big) integer expressed as a string
 * @return double
 */
double log_e(char* input_int);
