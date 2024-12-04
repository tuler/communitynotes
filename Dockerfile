FROM python:3.10-slim AS base
#FROM --platform=linux/riscv64 cartesi/python:3.10-jammy AS base-riscv64

RUN <<EOF
apt-get update
apt-get install -y --no-install-recommends \
    build-essential \
    cmake \
    libopenblas0 \
    ninja-build
EOF

WORKDIR /workspace
COPY requirements.txt .

#RUN pip config set global.extra-index-url https://tuler.github.io/riscv-python-wheels
RUN pip install -r requirements.txt

COPY . .

WORKDIR /workspace/sourcecode
ENTRYPOINT [ "python", "main.py" ]
