#include "mpfr_calculation.h"

// Guide: https://www.mpfr.org/mpfr-current/mpfr.html#index-mpfr_005frnd_005ft

double log_e(char* input_int) {
    mpfr_t loaded_f, result_f;

    mpfr_init2(loaded_f, 200);
    mpfr_set_str(loaded_f, input_int, 10, MPFR_RNDD);

    mpfr_init2(result_f, 200);


    mpfr_log(result_f, loaded_f, MPFR_RNDD);

    double result = mpfr_get_d(result_f, MPFR_RNDD);

    // mpfr_set_str(loaded_f, input_int, 10, )
    mpfr_clear(loaded_f);
    mpfr_clear(result_f);
    mpfr_free_cache();

    return result;
}

// #cgo LDFLAGS: -L "/opt/homebrew/lib" -lmpfr -lgmp
// #cgo CPPFLAGS: -I "/opt/homebrew/include"
// #cgo FLAGS: -fuse-ld=lld
// 위의 것들 다 필요없이 cpp -> c 확장자 변경으로 linker 관련 문제 해결
