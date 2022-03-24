#include "mpfr_calculation.h"

double log_e(char* input_int) {
    mpfr_t loaded_f, result_f;

    mpfr_init2(loaded_f, 200);
    mpfr_set_str(loaded_f, input_int, 10, MPFR_RNDD);

    mpfr_init2(result_f, 200);


    mpfr_log(result_f, loaded_f, MPFR_RNDD);

    double result = mpfr_get_d(result_f, MPFR_RNDD);

    mpfr_clear(loaded_f);
    mpfr_clear(result_f);
    mpfr_free_cache();

    return result;
}
