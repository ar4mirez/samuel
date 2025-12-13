# CUDA Guide

> **Applies to**: CUDA 11+, GPU Computing, Deep Learning, Scientific Computing

---

## Core Principles

1. **Parallelism First**: Design for thousands of concurrent threads
2. **Memory Hierarchy**: Optimize for global, shared, and register memory
3. **Coalesced Access**: Align memory access patterns for bandwidth
4. **Occupancy**: Maximize GPU utilization
5. **Host-Device Coordination**: Minimize data transfers

---

## Language-Specific Guardrails

### CUDA Version & Setup
- ✓ Use CUDA 11.0+ (12.x recommended for new projects)
- ✓ Use CMake with CUDA language support
- ✓ Target appropriate compute capability for hardware
- ✓ Use CUDA Toolkit matching driver version
- ✓ Enable `-arch=native` or specific `-arch=sm_XX` for optimization

### Code Style
- ✓ Use `__global__` for kernel functions
- ✓ Use `__device__` for device-only functions
- ✓ Use `__host__ __device__` for functions callable from both
- ✓ Prefix kernel names descriptively: `kernel_`, `k_`
- ✓ Use `snake_case` for variables and functions
- ✓ Use `PascalCase` for classes and structs
- ✓ Document thread/block dimensions in kernel comments

### Memory Management
- ✓ Use CUDA Unified Memory (`cudaMallocManaged`) for simplicity
- ✓ Use explicit `cudaMalloc`/`cudaMemcpy` for performance-critical code
- ✓ Always check CUDA API return values
- ✓ Use `cudaFree` for every `cudaMalloc`
- ✓ Prefer `cudaMemcpyAsync` with streams for overlap
- ✓ Use pinned memory (`cudaMallocHost`) for faster transfers

### Error Handling
- ✓ Check every CUDA API call with `cudaGetLastError()`
- ✓ Use error checking macro for all CUDA calls
- ✓ Handle out-of-memory errors gracefully
- ✓ Synchronize before checking kernel errors
- ✓ Use `cuda-memcheck` during development

### Thread Organization
- ✓ Choose block size as multiple of warp size (32)
- ✓ Common block sizes: 128, 256, 512 threads
- ✓ Ensure thread count covers all data elements
- ✓ Avoid divergent branches within warps
- ✓ Use grid-stride loops for flexible sizing

---

## Project Structure

### CMake CUDA Project
```
myproject/
├── CMakeLists.txt
├── include/
│   └── myproject/
│       ├── kernels.cuh           # Kernel declarations
│       └── cuda_utils.cuh        # Utilities
├── src/
│   ├── main.cpp                  # Host code
│   ├── kernels.cu                # Kernel implementations
│   └── cuda_utils.cu
├── tests/
│   └── test_kernels.cu
└── README.md
```

### CMakeLists.txt
```cmake
cmake_minimum_required(VERSION 3.18)
project(myproject LANGUAGES CXX CUDA)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CUDA_STANDARD 17)

# Find CUDA
find_package(CUDAToolkit REQUIRED)

# Set architecture (adjust for your GPU)
set(CMAKE_CUDA_ARCHITECTURES 70 75 80 86)

# Enable separable compilation for device linking
set(CMAKE_CUDA_SEPARABLE_COMPILATION ON)

# Library
add_library(myproject
    src/kernels.cu
    src/cuda_utils.cu
)

target_include_directories(myproject PUBLIC include)

# Executable
add_executable(main src/main.cpp)
target_link_libraries(main myproject CUDA::cudart)

# Tests
enable_testing()
add_executable(tests tests/test_kernels.cu)
target_link_libraries(tests myproject CUDA::cudart)
```

---

## Basic CUDA Patterns

### Error Checking Macro
```cuda
#include <cuda_runtime.h>
#include <cstdio>

#define CUDA_CHECK(call)                                                    \
    do {                                                                    \
        cudaError_t error = call;                                          \
        if (error != cudaSuccess) {                                        \
            fprintf(stderr, "CUDA error at %s:%d: %s\n",                   \
                    __FILE__, __LINE__, cudaGetErrorString(error));        \
            exit(EXIT_FAILURE);                                            \
        }                                                                   \
    } while (0)

#define CUDA_CHECK_LAST()                                                   \
    do {                                                                    \
        cudaError_t error = cudaGetLastError();                            \
        if (error != cudaSuccess) {                                        \
            fprintf(stderr, "CUDA kernel error at %s:%d: %s\n",            \
                    __FILE__, __LINE__, cudaGetErrorString(error));        \
            exit(EXIT_FAILURE);                                            \
        }                                                                   \
    } while (0)
```

### Simple Kernel
```cuda
// Vector addition kernel
__global__ void vector_add(const float* a, const float* b, float* c, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;

    if (idx < n) {
        c[idx] = a[idx] + b[idx];
    }
}

// Launch kernel
void launch_vector_add(const float* a, const float* b, float* c, int n) {
    int block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;

    vector_add<<<grid_size, block_size>>>(a, b, c, n);
    CUDA_CHECK_LAST();
}
```

### Grid-Stride Loop Pattern
```cuda
// Handles any size array efficiently
__global__ void vector_add_stride(const float* a, const float* b, float* c, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    int stride = blockDim.x * gridDim.x;

    for (int i = idx; i < n; i += stride) {
        c[i] = a[i] + b[i];
    }
}

// Can use fixed grid size
void launch_vector_add_stride(const float* a, const float* b, float* c, int n) {
    int block_size = 256;
    int grid_size = 256;  // Fixed, handles any n

    vector_add_stride<<<grid_size, block_size>>>(a, b, c, n);
    CUDA_CHECK_LAST();
}
```

---

## Memory Management

### Explicit Memory Management
```cuda
#include <cuda_runtime.h>

void example_explicit_memory() {
    const int n = 1000000;
    const size_t size = n * sizeof(float);

    // Host memory
    float* h_a = new float[n];
    float* h_b = new float[n];
    float* h_c = new float[n];

    // Initialize host data
    for (int i = 0; i < n; i++) {
        h_a[i] = 1.0f;
        h_b[i] = 2.0f;
    }

    // Device memory
    float *d_a, *d_b, *d_c;
    CUDA_CHECK(cudaMalloc(&d_a, size));
    CUDA_CHECK(cudaMalloc(&d_b, size));
    CUDA_CHECK(cudaMalloc(&d_c, size));

    // Copy to device
    CUDA_CHECK(cudaMemcpy(d_a, h_a, size, cudaMemcpyHostToDevice));
    CUDA_CHECK(cudaMemcpy(d_b, h_b, size, cudaMemcpyHostToDevice));

    // Launch kernel
    int block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    vector_add<<<grid_size, block_size>>>(d_a, d_b, d_c, n);
    CUDA_CHECK_LAST();

    // Copy back
    CUDA_CHECK(cudaMemcpy(h_c, d_c, size, cudaMemcpyDeviceToHost));

    // Cleanup
    CUDA_CHECK(cudaFree(d_a));
    CUDA_CHECK(cudaFree(d_b));
    CUDA_CHECK(cudaFree(d_c));
    delete[] h_a;
    delete[] h_b;
    delete[] h_c;
}
```

### Unified Memory
```cuda
void example_unified_memory() {
    const int n = 1000000;
    const size_t size = n * sizeof(float);

    // Unified memory (accessible from host and device)
    float *a, *b, *c;
    CUDA_CHECK(cudaMallocManaged(&a, size));
    CUDA_CHECK(cudaMallocManaged(&b, size));
    CUDA_CHECK(cudaMallocManaged(&c, size));

    // Initialize on host (no explicit copy needed)
    for (int i = 0; i < n; i++) {
        a[i] = 1.0f;
        b[i] = 2.0f;
    }

    // Launch kernel
    int block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    vector_add<<<grid_size, block_size>>>(a, b, c, n);

    // Wait for GPU to finish
    CUDA_CHECK(cudaDeviceSynchronize());

    // Access result on host (no explicit copy)
    printf("Result: %f\n", c[0]);

    // Cleanup
    CUDA_CHECK(cudaFree(a));
    CUDA_CHECK(cudaFree(b));
    CUDA_CHECK(cudaFree(c));
}
```

### Pinned Memory for Fast Transfers
```cuda
void example_pinned_memory() {
    const int n = 1000000;
    const size_t size = n * sizeof(float);

    // Pinned (page-locked) host memory
    float* h_a;
    CUDA_CHECK(cudaMallocHost(&h_a, size));

    float* d_a;
    CUDA_CHECK(cudaMalloc(&d_a, size));

    // Faster transfer with pinned memory
    CUDA_CHECK(cudaMemcpy(d_a, h_a, size, cudaMemcpyHostToDevice));

    // Cleanup
    CUDA_CHECK(cudaFree(d_a));
    CUDA_CHECK(cudaFreeHost(h_a));
}
```

---

## Shared Memory

### Basic Shared Memory Usage
```cuda
__global__ void reduce_sum(const float* input, float* output, int n) {
    extern __shared__ float sdata[];

    int tid = threadIdx.x;
    int idx = blockIdx.x * blockDim.x + threadIdx.x;

    // Load to shared memory
    sdata[tid] = (idx < n) ? input[idx] : 0.0f;
    __syncthreads();

    // Reduction in shared memory
    for (int s = blockDim.x / 2; s > 0; s >>= 1) {
        if (tid < s) {
            sdata[tid] += sdata[tid + s];
        }
        __syncthreads();
    }

    // Write result
    if (tid == 0) {
        output[blockIdx.x] = sdata[0];
    }
}

// Launch with dynamic shared memory
void launch_reduce(const float* input, float* output, int n) {
    int block_size = 256;
    int grid_size = (n + block_size - 1) / block_size;
    size_t shared_size = block_size * sizeof(float);

    reduce_sum<<<grid_size, block_size, shared_size>>>(input, output, n);
    CUDA_CHECK_LAST();
}
```

### Tiled Matrix Multiplication
```cuda
#define TILE_SIZE 16

__global__ void matrix_multiply_tiled(
    const float* A, const float* B, float* C,
    int M, int N, int K
) {
    __shared__ float As[TILE_SIZE][TILE_SIZE];
    __shared__ float Bs[TILE_SIZE][TILE_SIZE];

    int row = blockIdx.y * TILE_SIZE + threadIdx.y;
    int col = blockIdx.x * TILE_SIZE + threadIdx.x;

    float sum = 0.0f;

    for (int t = 0; t < (K + TILE_SIZE - 1) / TILE_SIZE; t++) {
        // Load tiles into shared memory
        int a_col = t * TILE_SIZE + threadIdx.x;
        int b_row = t * TILE_SIZE + threadIdx.y;

        As[threadIdx.y][threadIdx.x] =
            (row < M && a_col < K) ? A[row * K + a_col] : 0.0f;
        Bs[threadIdx.y][threadIdx.x] =
            (b_row < K && col < N) ? B[b_row * N + col] : 0.0f;

        __syncthreads();

        // Compute partial dot product
        for (int k = 0; k < TILE_SIZE; k++) {
            sum += As[threadIdx.y][k] * Bs[k][threadIdx.x];
        }

        __syncthreads();
    }

    if (row < M && col < N) {
        C[row * N + col] = sum;
    }
}
```

---

## Streams and Async Operations

### Concurrent Kernels with Streams
```cuda
void example_streams() {
    const int n = 1000000;
    const size_t size = n * sizeof(float);
    const int num_streams = 4;

    // Create streams
    cudaStream_t streams[num_streams];
    for (int i = 0; i < num_streams; i++) {
        CUDA_CHECK(cudaStreamCreate(&streams[i]));
    }

    // Allocate memory
    float *h_a, *d_a;
    CUDA_CHECK(cudaMallocHost(&h_a, size * num_streams));
    CUDA_CHECK(cudaMalloc(&d_a, size * num_streams));

    // Launch operations in parallel streams
    for (int i = 0; i < num_streams; i++) {
        size_t offset = i * n;

        // Async copy
        CUDA_CHECK(cudaMemcpyAsync(
            d_a + offset, h_a + offset, size,
            cudaMemcpyHostToDevice, streams[i]
        ));

        // Kernel in stream
        int block_size = 256;
        int grid_size = (n + block_size - 1) / block_size;
        some_kernel<<<grid_size, block_size, 0, streams[i]>>>(
            d_a + offset, n
        );

        // Async copy back
        CUDA_CHECK(cudaMemcpyAsync(
            h_a + offset, d_a + offset, size,
            cudaMemcpyDeviceToHost, streams[i]
        ));
    }

    // Wait for all streams
    CUDA_CHECK(cudaDeviceSynchronize());

    // Cleanup
    for (int i = 0; i < num_streams; i++) {
        CUDA_CHECK(cudaStreamDestroy(streams[i]));
    }
    CUDA_CHECK(cudaFreeHost(h_a));
    CUDA_CHECK(cudaFree(d_a));
}
```

### Events for Timing
```cuda
void benchmark_kernel() {
    cudaEvent_t start, stop;
    CUDA_CHECK(cudaEventCreate(&start));
    CUDA_CHECK(cudaEventCreate(&stop));

    // Record start
    CUDA_CHECK(cudaEventRecord(start));

    // Launch kernel
    some_kernel<<<grid_size, block_size>>>(args);

    // Record stop
    CUDA_CHECK(cudaEventRecord(stop));
    CUDA_CHECK(cudaEventSynchronize(stop));

    // Calculate elapsed time
    float milliseconds = 0;
    CUDA_CHECK(cudaEventElapsedTime(&milliseconds, start, stop));
    printf("Kernel time: %.3f ms\n", milliseconds);

    CUDA_CHECK(cudaEventDestroy(start));
    CUDA_CHECK(cudaEventDestroy(stop));
}
```

---

## Atomic Operations

### Common Atomics
```cuda
__global__ void histogram(const int* data, int* hist, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;

    if (idx < n) {
        int bin = data[idx];
        atomicAdd(&hist[bin], 1);
    }
}

__global__ void find_max(const float* data, float* result, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;

    if (idx < n) {
        // atomicMax for floats (CUDA 11+)
        atomicMax((int*)result, __float_as_int(data[idx]));
    }
}

__global__ void cas_example(int* data) {
    int idx = threadIdx.x;

    // Compare-and-swap
    int old = atomicCAS(&data[0], 0, idx);
    // Only one thread succeeds in setting value
}
```

### Shared Memory Atomics (Faster)
```cuda
__global__ void histogram_shared(const int* data, int* hist, int n, int num_bins) {
    extern __shared__ int s_hist[];

    int tid = threadIdx.x;
    int idx = blockIdx.x * blockDim.x + threadIdx.x;

    // Initialize shared histogram
    for (int i = tid; i < num_bins; i += blockDim.x) {
        s_hist[i] = 0;
    }
    __syncthreads();

    // Accumulate in shared memory
    if (idx < n) {
        atomicAdd(&s_hist[data[idx]], 1);
    }
    __syncthreads();

    // Write to global memory
    for (int i = tid; i < num_bins; i += blockDim.x) {
        atomicAdd(&hist[i], s_hist[i]);
    }
}
```

---

## Thrust Library

### Using Thrust for High-Level Operations
```cuda
#include <thrust/device_vector.h>
#include <thrust/host_vector.h>
#include <thrust/sort.h>
#include <thrust/reduce.h>
#include <thrust/transform.h>
#include <thrust/functional.h>

void example_thrust() {
    // Host vector
    thrust::host_vector<float> h_vec(1000000);
    for (int i = 0; i < h_vec.size(); i++) {
        h_vec[i] = static_cast<float>(rand()) / RAND_MAX;
    }

    // Copy to device
    thrust::device_vector<float> d_vec = h_vec;

    // Sort
    thrust::sort(d_vec.begin(), d_vec.end());

    // Reduce (sum)
    float sum = thrust::reduce(d_vec.begin(), d_vec.end(), 0.0f, thrust::plus<float>());
    printf("Sum: %f\n", sum);

    // Transform
    thrust::transform(
        d_vec.begin(), d_vec.end(),
        d_vec.begin(),
        thrust::negate<float>()
    );

    // Custom functor
    struct square_functor {
        __host__ __device__
        float operator()(float x) const {
            return x * x;
        }
    };

    thrust::transform(
        d_vec.begin(), d_vec.end(),
        d_vec.begin(),
        square_functor()
    );

    // Copy back
    h_vec = d_vec;
}
```

---

## Performance Optimization

### Optimization Guardrails
- ✓ Maximize occupancy (threads per SM)
- ✓ Coalesce global memory accesses
- ✓ Use shared memory for reused data
- ✓ Avoid warp divergence
- ✓ Minimize host-device transfers
- ✓ Use streams for overlap
- ✓ Profile with Nsight Compute

### Memory Coalescing
```cuda
// BAD: Strided access (not coalesced)
__global__ void bad_access(float* data, int stride, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        data[idx * stride] = 1.0f;  // Non-coalesced!
    }
}

// GOOD: Contiguous access (coalesced)
__global__ void good_access(float* data, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        data[idx] = 1.0f;  // Coalesced
    }
}

// For 2D arrays: access by row (contiguous)
__global__ void row_access(float* matrix, int rows, int cols) {
    int row = blockIdx.y * blockDim.y + threadIdx.y;
    int col = blockIdx.x * blockDim.x + threadIdx.x;

    if (row < rows && col < cols) {
        // Row-major: contiguous within warp if threads vary in col
        matrix[row * cols + col] = 1.0f;
    }
}
```

### Avoiding Warp Divergence
```cuda
// BAD: Divergent branch
__global__ void divergent_kernel(float* data, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        if (idx % 2 == 0) {
            data[idx] = expensive_operation_a(idx);
        } else {
            data[idx] = expensive_operation_b(idx);
        }
    }
}

// BETTER: Process separately or use predication
__global__ void less_divergent_kernel(float* data, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        // Both computed, compiler may optimize
        float a = expensive_operation_a(idx);
        float b = expensive_operation_b(idx);
        data[idx] = (idx % 2 == 0) ? a : b;
    }
}
```

---

## Testing

### Simple Test Framework
```cuda
#include <cuda_runtime.h>
#include <cstdio>
#include <cmath>

#define ASSERT_NEAR(expected, actual, epsilon)                    \
    do {                                                          \
        if (fabs((expected) - (actual)) > (epsilon)) {           \
            fprintf(stderr, "FAILED: %s:%d: expected %f, got %f\n", \
                    __FILE__, __LINE__, (expected), (actual));    \
            return false;                                         \
        }                                                         \
    } while (0)

bool test_vector_add() {
    const int n = 1000;
    float *a, *b, *c;

    CUDA_CHECK(cudaMallocManaged(&a, n * sizeof(float)));
    CUDA_CHECK(cudaMallocManaged(&b, n * sizeof(float)));
    CUDA_CHECK(cudaMallocManaged(&c, n * sizeof(float)));

    // Initialize
    for (int i = 0; i < n; i++) {
        a[i] = 1.0f;
        b[i] = 2.0f;
    }

    // Run kernel
    launch_vector_add(a, b, c, n);
    CUDA_CHECK(cudaDeviceSynchronize());

    // Verify
    for (int i = 0; i < n; i++) {
        ASSERT_NEAR(3.0f, c[i], 1e-5f);
    }

    CUDA_CHECK(cudaFree(a));
    CUDA_CHECK(cudaFree(b));
    CUDA_CHECK(cudaFree(c));

    printf("test_vector_add PASSED\n");
    return true;
}

int main() {
    bool all_passed = true;

    all_passed &= test_vector_add();
    // Add more tests...

    return all_passed ? 0 : 1;
}
```

---

## Common Pitfalls

### Don't Do This
```cuda
// Forgetting bounds check
__global__ void bad_bounds(float* data, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    data[idx] = 1.0f;  // Out of bounds if idx >= n!
}

// Not synchronizing before checking errors
kernel<<<grid, block>>>(args);
// cudaGetLastError() here won't catch async kernel errors

// Memory leak
float* d_ptr;
cudaMalloc(&d_ptr, size);
// Forgot cudaFree(d_ptr)!

// Race condition in shared memory
__global__ void race_condition() {
    __shared__ float shared_data[256];
    shared_data[threadIdx.x] = threadIdx.x;
    // Missing __syncthreads()!
    float val = shared_data[threadIdx.x + 1];  // Race!
}
```

### Do This Instead
```cuda
// Always check bounds
__global__ void good_bounds(float* data, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;
    if (idx < n) {
        data[idx] = 1.0f;
    }
}

// Synchronize before checking errors
kernel<<<grid, block>>>(args);
CUDA_CHECK(cudaDeviceSynchronize());
CUDA_CHECK_LAST();

// RAII wrapper for CUDA memory
template<typename T>
class DeviceBuffer {
public:
    explicit DeviceBuffer(size_t count) : size_(count * sizeof(T)) {
        CUDA_CHECK(cudaMalloc(&ptr_, size_));
    }

    ~DeviceBuffer() {
        cudaFree(ptr_);
    }

    T* get() { return ptr_; }

    // Non-copyable
    DeviceBuffer(const DeviceBuffer&) = delete;
    DeviceBuffer& operator=(const DeviceBuffer&) = delete;

private:
    T* ptr_;
    size_t size_;
};

// Proper synchronization
__global__ void no_race() {
    __shared__ float shared_data[256];
    shared_data[threadIdx.x] = threadIdx.x;
    __syncthreads();  // Wait for all threads
    float val = shared_data[(threadIdx.x + 1) % blockDim.x];
}
```

---

## Debugging Tools

### cuda-memcheck
```bash
# Check for memory errors
cuda-memcheck ./program

# Check for race conditions
cuda-memcheck --tool racecheck ./program

# Check for initialization errors
cuda-memcheck --tool initcheck ./program

# Memory leak detection
cuda-memcheck --leak-check full ./program
```

### Nsight Compute (Profiling)
```bash
# Profile kernel
ncu ./program

# Profile specific kernel
ncu --kernel-name "my_kernel" ./program

# Generate report
ncu -o report ./program
ncu -i report.ncu-rep  # View report
```

### Device Query
```cuda
void print_device_info() {
    int device;
    cudaGetDevice(&device);

    cudaDeviceProp props;
    cudaGetDeviceProperties(&props, device);

    printf("Device: %s\n", props.name);
    printf("Compute capability: %d.%d\n", props.major, props.minor);
    printf("Total global memory: %.2f GB\n",
           props.totalGlobalMem / (1024.0 * 1024.0 * 1024.0));
    printf("Shared memory per block: %zu KB\n",
           props.sharedMemPerBlock / 1024);
    printf("Max threads per block: %d\n", props.maxThreadsPerBlock);
    printf("Warp size: %d\n", props.warpSize);
    printf("Max grid size: (%d, %d, %d)\n",
           props.maxGridSize[0], props.maxGridSize[1], props.maxGridSize[2]);
}
```

---

## References

- [CUDA C++ Programming Guide](https://docs.nvidia.com/cuda/cuda-c-programming-guide/)
- [CUDA C++ Best Practices Guide](https://docs.nvidia.com/cuda/cuda-c-best-practices-guide/)
- [Thrust Documentation](https://nvidia.github.io/thrust/)
- [Nsight Compute Documentation](https://docs.nvidia.com/nsight-compute/)
- [CUDA Samples](https://github.com/NVIDIA/cuda-samples)
- [CUDA Toolkit Documentation](https://docs.nvidia.com/cuda/)
