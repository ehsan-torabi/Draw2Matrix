# Draw2Matrix: Hand-Drawn Pattern to Matrix Converter

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/dl/)
[![Fyne Version](https://img.shields.io/badge/Fyne-v2.x-7F4FC5?style=flat)](https://fyne.io/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-blue)](https://github.com/ehsan-torabi/Draw2Matrix/releases/latest)

> Convert hand-drawn patterns into machine learning-ready matrices with real-time processing and multiple export formats.

Draw2Matrix is a powerful Go-based drawing application that transforms hand-drawn patterns into machine-readable binary matrices. Built with the Fyne toolkit, it offers real-time conversion of drawings into various formats including CSV, MATLAB (with One-Hot Encoding), and PNG. Perfect for:

- 🎓 Creating machine learning datasets
- 📊 Pattern recognition research
- 🤖 AI/ML training data preparation
- 🎨 Digital pattern visualization
- 📚 Educational tools development

<div align="center">
  <img src="Icon.png" alt="Draw2Matrix Application Interface - Pattern to Matrix Converter" width="64" height="64">
</div>

## ✨ What's New (August 2025)

- **Enhanced Machine Learning Support**:

  - One-Hot Encoding for MATLAB export - Perfect for neural networks
  - Matrix counter with progress tracking
  - Real-time status updates with animations

- **Improved Drawing Tools**:

  - Responsive paint window with precise input
  - Advanced image processing for accurate pattern recognition
  - Real-time preview of matrix conversion

- **Performance Optimizations**:
  - Efficient light theme UI for better pattern visibility
  - Streamlined initialization process
  - Optimized codebase for faster processing

## 🚀 Key Features

### Pattern Recognition & Conversion

- Real-time conversion of drawings to binary matrices
- Automated pattern detection and processing
- Intelligent matrix size adaptation

### Export Flexibility

- Multiple format support:
  - CSV export with optional flattening
  - MATLAB format with One-Hot Encoding
  - High-resolution PNG image export
- Batch processing capabilities
- Custom label support

### User Experience

- Intuitive drawing interface
- Dynamic matrix size adjustment
- Customizable output options:
  - Row/Column flattening
  - Matrix dimension control
  - Label management system

## 🔧 Installation

### System Requirements

- **Go**: Version 1.22 or later
- **Fyne**: v2.x toolkit
- **OS Support**: Windows, macOS, Linux

### Quick Start

```bash
# Clone the repository
git clone https://github.com/ehsan-torabi/Draw2Matrix.git

# Navigate to project directory
cd Draw2Matrix

# Install dependencies
go mod download

# Run the application
go run .
```

## 📝 Usage Guide

1. **Starting Up**:

   ```bash
   # Run the compiled executable
   ./Draw2Matrix

   # Or use Go directly
   go run .
   ```

2. **Matrix Configuration**:

   - Set your desired matrix dimensions
   - Choose output format:
     - Standard matrix
     - Flattened matrix
     - One-Hot encoded (MATLAB)
   - Apply settings with "Save Settings"

3. **Drawing Interface**:

   - Use the enhanced paint window for drawing
   - Real-time matrix conversion
   - Track additions with the matrix counter
   - Clear canvas option available

4. **Export Process**:
   - Add descriptive labels
   - Select export directory
   - Choose format:
     - CSV (with flattening options)
     - MATLAB (with One-Hot encoding)
     - PNG image
   - Monitor progress through animated status updates

## 📊 Output Formats

### CSV Export

```csv
Input,Target
[1 0 1 0],label      # Flattened format
```

### MATLAB Export

The application generates optimized MATLAB-compatible files:

1. **Matrix Data** (`data.txt`):

   ```matlab
   % Standard Format
   [ 1 0 1;
     0 1 0;
     1 1 0 ]

   ```

2. **Labels** (`target.txt`):

   ```matlab
   [ 'A' 'B' 'C' ]

   % One-Hot Encoded Format
   [ 1 0 0;
     0 1 0;
     0 0 1 ]
   ```

## 🗂️ Project Structure

- **Core Components**:
  - `main.go`: Entry point and core application logic
  - `paintWindow.go`: Enhanced drawing interface
  - `imageTools.go`: Advanced image processing
  - `dataTools.go`: Data handling and export functions
  - `controlFunctions.go`: UI control management
  - `customWidget.go`: Custom widget implementations

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🌟 Community & Support

### Contributing

We welcome contributions! Here's how you can help:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Technology Stack

- [Go](https://go.dev/) - Modern, fast programming language
- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [bild](https://github.com/anthonynsimon/bild) - Advanced image processing

### Connect & Support

- 📧 **Author**: Ehsan Torabi
- 💬 **Telegram**: [@ehsan_torabi_frs](https://t.me/ehsan_torabi_frs)
- 🌟 **Project**: [Draw2Matrix on GitHub](https://github.com/ehsan-torabi/Draw2Matrix)
- 📄 **License**: MIT License - [View License](LICENSE)

### Keywords

`pattern recognition`, `machine learning`, `data preprocessing`, `matrix conversion`, `golang`, `fyne`, `drawing tool`, `binary matrix`, `dataset creation`, `educational tool`, `AI training`, `cross-platform`, `open source`
