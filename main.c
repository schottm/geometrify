#include "org_iMage_geometrify_ParallelTrianglePictureFilter.h"
#include "main.h"
#include <stdio.h>
#include <stdlib.h>

JNIEXPORT jintArray JNICALL Java_org_iMage_geometrify_ParallelTrianglePictureFilter_calculate
    (JNIEnv * env, jclass clazz, jint width, jint height, jintArray data, jint iterations, jint samples, jint threads, jlong seed) {

    //convert jintArray to c_data_ptr
    jint *c_data_ptr;
    c_data_ptr = (*env)->GetIntArrayElements(env, data, NULL);

    //create go_data_ptr
    GoInt32 *go_data_ptr = 0;
    go_data_ptr = (GoInt32*) malloc(width * height * sizeof(GoInt32));

    //copy c_data_ptr to go_data_ptr
    for(int i = 0; i < width * height; ++i) {
        *(go_data_ptr + i) = (GoInt32) *(c_data_ptr + i);
    }

    GoInt go_width = (GoInt) width;
    GoInt go_height = (GoInt) height;

    GoSlice go_data_slice = {go_data_ptr, width * height, width * height};

    GoInt go_iterations = (GoInt) iterations;
    GoInt go_samples = (GoInt) samples;
    GoInt go_threads = (GoInt) threads;

    GoInt64 go_seed = (GoInt64) seed;

    GoSlice go_result = Calculate(go_width, go_height, go_data_slice, go_iterations, go_samples, go_threads, go_seed);

    for (int j = 0; j < width * height; j++) {
       *(c_data_ptr + j) = (jint) ((GoInt32 *)go_result.data)[j];
    }

    //convert jint[] to jintArray
    jintArray out_ints;
    out_ints = (*env)->NewIntArray(env, width * height);
    (*env)->SetIntArrayRegion(env, out_ints, 0, width * height, c_data_ptr);

    return out_ints;
}