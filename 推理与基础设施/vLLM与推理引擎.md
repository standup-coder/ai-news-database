---
线索: vLLM与推理引擎
主题: 推理与基础设施
别名: [vLLM, SGLang, TensorRT-LLM, 推理优化]
状态: 活跃
创建: 2026-07-30
更新: 2026-08-03
关键角色: [vLLM 社区, UC Berkeley, NVIDIA, LMSYS, DeepSeek]
---

# vLLM与推理引擎

> 开源推理引擎的军备竞赛：PagedAttention 开启的高吞吐推理时代，决定每个 token 的成本下限。

## 概述

vLLM 起源于 UC Berkeley 的 PagedAttention 论文（2023），以类操作系统虚拟内存的方式管理 KV Cache，大幅提升吞吐量，迅速成为开源 LLM 推理的事实标准。其竞争者包括 SGLang（RadixAttention）、NVIDIA TensorRT-LLM、Hugging Face TGI 等。推理引擎的效率直接决定 API 价格战的成本底线，是「模型民主化」背后最重要的基础设施。

该线索追踪推理引擎的技术演进（连续批处理、投机解码、前缀缓存、分布式推理）与生态格局。

## 时间线

### 2026-04

- **2026-04-24** · [vLLM 整合 DeepSeek V4 优化](https://vllm.ai/blog/2026-04-24-deepseek-v4)
  vLLM 整合 FlashMLA + FlashInfer，为 DeepSeek V4 提供长上下文推理优化（DeepSeek Sparse Attention 有界注意力 + MoE 内核）。

### 2026（持续）

- **2026（全年）** · vLLM / SGLang / TensorRT-LLM 三足鼎立
  TensorRT-LLM 在 NVIDIA 硬件上原始吞吐最高；SGLang 在结构化输出/受限解码占优；vLLM 易用、模型覆盖广，均已生产成熟。
- **2026（持续）** · DeepSeek FlashMLA / DSA
  FlashMLA CUDA 内核：预填充 640 TFlops、解码 410 TFlops；DeepSeek Sparse Attention 约降低 50% 推理成本。

### 2025-05

- **2025-05** · 大规模 MoE 推理成为主战场
  DeepSeek-V3/R1 等巨型 MoE 模型的流行推动 vLLM/SGLang 竞争 PD 分离（Prefill/Decode 解耦）、专家并行与多机推理能力，llm-d、Dynamo 等编排层项目涌现。

### 2025-01

- **2025-01** · vLLM V1 架构重构发布
  核心引擎重写：更干净的调度器、原生支持分块预填充与前缀缓存，吞吐与延迟全面提升，巩固其对 TGI 等早期方案的替代。

### 2024-07

- **2024-07** · vLLM 成为 Llama 3.1 405B 官方推荐推理方案之一
  Meta 发布 405B 模型时与 vLLM 社区合作支持多机推理，标志开源推理引擎具备服务最大规模开源模型的能力。

### 2023-06

- **2023-06** · [vLLM 开源发布](https://blog.vllm.ai/2023/06/20/vllm.html)
  PagedAttention 论文配套开源实现，宣称相对 HF Transformers 提升最高 24 倍吞吐，LMSYS 用其支撑 Vicuna 演示服务，项目迅速爆红。

## 分析

1. **推理成本是产业杠杆**：过去两年 LLM API 价格下降了两个数量级，推理引擎优化（批处理、量化、投机解码）贡献了其中相当大比例。每一次引擎效率突破都直接转化为应用层的可行性边界扩张。

2. **开源治理的样本**：vLLM 从伯克利项目成长为 PyTorch 基金会旗下的中立项目，贡献者覆盖所有芯片厂商（NVIDIA/AMD/Intel/TPU/昇腾）——它已成为硬件竞争的「中立缓冲层」，谁都输不起。

3. **从单机到集群**：推理优化的前沿已从单卡 kernel 转向集群级调度（PD 分离、KV 传输、路由），推理系统正在复刻当年大数据系统（Hadoop→Spark→K8s）的演化轨迹。

## 关联线索

- [[推理与基础设施/AI芯片竞争]]
- [[开源模型与生态/Llama]]
- [[开源模型与生态/DeepSeek]]
