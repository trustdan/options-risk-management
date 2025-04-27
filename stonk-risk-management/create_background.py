#!/usr/bin/env python3
"""
Create a DMG background image with installation instructions
"""
try:
    from PIL import Image, ImageDraw, ImageFont
except ImportError:
    print("Pillow library is required. Please install it with: pip install Pillow")
    exit(1)

# Create image with dimensions matching DMG window size
width = 800
height = 400
background_color = (240, 240, 240)  # Light gray

# Create a new image with the specified dimensions and color
img = Image.new('RGB', (width, height), background_color)
draw = ImageDraw.Draw(img)

# Draw a light gradient background
for y in range(height):
    # Create a subtle gradient from top to bottom
    color = (240, 240, 240 + int(y * 15 / height))
    draw.line([(0, y), (width, y)], fill=color)

# Try to find a font that's likely to be on macOS
try:
    # Try SF Pro first (macOS system font)
    font_large = ImageFont.truetype("/System/Library/Fonts/SFNSDisplay.ttf", 32)
    font_small = ImageFont.truetype("/System/Library/Fonts/SFNSDisplay.ttf", 20)
except:
    try:
        # Try SF Pro (alternate location)
        font_large = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Display-Regular.otf", 32)
        font_small = ImageFont.truetype("/System/Library/Fonts/SF-Pro-Display-Regular.otf", 20)
    except:
        try:
            # Try Helvetica as fallback
            font_large = ImageFont.truetype("/System/Library/Fonts/Helvetica.ttc", 32)
            font_small = ImageFont.truetype("/System/Library/Fonts/Helvetica.ttc", 20)
        except:
            # If specific fonts can't be found, use default
            font_large = ImageFont.load_default()
            font_small = ImageFont.load_default()

# Add title
title = "Options Trading Dashboard"
draw.text((width//2, 60), title, fill=(0, 150, 0), font=font_large, anchor="mm")

# Add instruction text
instructions = ["To install:", 
                "1. Drag the app to the Applications folder",
                "2. Eject the disk image",
                "3. Launch from Applications"]

# Position instructions in the lower part of the image
y_position = 220
for instruction in instructions:
    draw.text((width//2, y_position), instruction, fill=(50, 50, 50), font=font_small, anchor="mm")
    y_position += 30

# Add arrows to indicate dragging direction
arrow_x1, arrow_y = 240, 180  # Start position
arrow_x2, arrow_y = 360, 180  # End position
arrow_color = (0, 150, 0)  # Green

# Draw arrow line
line_width = 3
draw.line([(arrow_x1, arrow_y), (arrow_x2, arrow_y)], fill=arrow_color, width=line_width)

# Draw arrow head
head_size = 10
draw.polygon([(arrow_x2, arrow_y), 
              (arrow_x2 - head_size, arrow_y - head_size//2), 
              (arrow_x2 - head_size, arrow_y + head_size//2)], 
              fill=arrow_color)

# Save the image
output_path = "build/resources/dmg-background.png"
img.save(output_path)
print(f"Background image created: {output_path}") 