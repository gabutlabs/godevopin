# ==============================================================================
# Variabel Lingkungan & Perintah Dasar
# ==============================================================================
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
# Dapatkan Arsitektur dan OS dari lingkungan saat ini
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS)

# Jalur (Paths)
BASE_PATH := $(shell pwd)
BUILD_PATH = $(BASE_PATH)/build
WEB_PATH=$(BASE_PATH)/frontend
ASSET_PATH=$(BASE_PATH)/cmd/devopin/web/dist
CORE_MAIN=$(BASE_PATH)/cmd/devopin/main.go
BUILD_NAME=godevopin
CORE_PATH=$(BASE_PATH)/cmd/devopin

# ==============================================================================
# PHONY TARGETS
# Menghindari bentrokan jika ada file dengan nama yang sama (mis. "all", "clean")
# ==============================================================================
.PHONY: all clean clean_assets build_frontend build_core_on_linux build_core_on_darwin build_all

# ==============================================================================
# 1. TARGET DEFAULT
# 'all' harus menjadi rule pertama agar make menjadikannya target default
# ==============================================================================
all: build_all

# ==============================================================================
# 2. TARGET BUILD KESELURUHAN (Menggantikan build_all Anda)
# ==============================================================================
build_all: clean_assets build_frontend build_core_on_linux # Menghapus build_agent_on_linux karena rule-nya tidak didefinisikan

# ==============================================================================
# 3. TARGET FRONTEND
# ==============================================================================
clean_assets:
	@echo "--- Membersihkan aset frontend yang sudah di-build..."
	rm -rf $(ASSET_PATH)

build_frontend:
	@echo "--- Membangun frontend (pnpm install & pnpm build)..."
	cd $(WEB_PATH) && pnpm install && pnpm build
	cp -r $(WEB_PATH)/dist $(CORE_PATH)/web/

# ==============================================================================
# 4. TARGET BACKEND/CORE
# ==============================================================================
# Build untuk OS/Arsitektur yang sedang berjalan (Linux jika di Linux, Darwin jika di Darwin, dst.)
build_core_on_linux:
	@echo "--- Membangun Core Go untuk $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(BUILD_NAME) $(CORE_MAIN)

# Target terpisah jika perlu kompilasi silang (cross-compile) ke Linux dari Mac/Windows
build_core_on_darwin:
	@echo "--- Cross-compiling Core Go ke Linux/AMD64..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -trimpath -ldflags '-s -w' -o $(BUILD_PATH)/$(BUILD_NAME)_linux_amd64 $(CORE_MAIN)

# ==============================================================================
# 5. TARGET CLEANING
# ==============================================================================
clean: clean_assets
	@echo "--- Membersihkan binary core di folder build..."
	$(GOCLEAN)
	rm -rf $(BUILD_PATH)/$(BUILD_NAME)*