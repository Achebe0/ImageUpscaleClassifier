#!/usr/bin/env python3
"""
Simple PyTorch image upscaler CLI
Takes an input image and upscales it using PyTorch
"""

import argparse
import sys
import torch
import cv2
import numpy as np
from pathlib import Path


def upscale_image(input_path: str, output_path: str, scale: int = 2) -> bool:
    """
    Upscale an image using PyTorch/OpenCV
    
    Args:
        input_path: Path to input image
        output_path: Path to save upscaled image
        scale: Upscaling factor (2x, 4x, etc)
    
    Returns:
        True if successful, False otherwise
    """
    try:
        # Read image
        img = cv2.imread(input_path, cv2.IMREAD_COLOR)
        if img is None:
            print(f"ERROR: Could not read image: {input_path}", file=sys.stderr)
            return False
        
        h, w = img.shape[:2]
        print(f"Input image size: {w}x{h}", file=sys.stderr)
        
        # Upscale using OpenCV bicubic interpolation
        # For production, replace with actual PyTorch model (ESRGAN, RealESRGAN, etc)
        new_h = h * scale
        new_w = w * scale
        upscaled = cv2.resize(img, (new_w, new_h), interpolation=cv2.INTER_CUBIC)
        
        print(f"Output image size: {new_w}x{new_h}", file=sys.stderr)
        
        # Save upscaled image
        if not cv2.imwrite(output_path, upscaled):
            print(f"ERROR: Could not write image: {output_path}", file=sys.stderr)
            return False
        
        print(f"Upscaled image saved to: {output_path}", file=sys.stderr)
        return True
    
    except Exception as e:
        print(f"ERROR: Upscaling failed: {str(e)}", file=sys.stderr)
        return False


def main():
    parser = argparse.ArgumentParser(
        description="Upscale images using PyTorch"
    )
    parser.add_argument(
        "--input", 
        required=True, 
        help="Path to input image"
    )
    parser.add_argument(
        "--output", 
        required=True, 
        help="Path to output image"
    )
    parser.add_argument(
        "--scale", 
        type=int, 
        default=2, 
        help="Upscaling factor (default: 2)"
    )
    
    args = parser.parse_args()
    
    # Validate inputs
    input_path = Path(args.input)
    if not input_path.exists():
        print(f"ERROR: Input file not found: {args.input}", file=sys.stderr)
        sys.exit(1)
    
    output_path = Path(args.output)
    if args.scale not in [2, 3, 4]:
        print(f"ERROR: Scale must be 2, 3, or 4, got {args.scale}", file=sys.stderr)
        sys.exit(1)
    
    # Upscale image
    success = upscale_image(str(input_path), str(output_path), args.scale)
    sys.exit(0 if success else 1)


if __name__ == "__main__":
    main()
