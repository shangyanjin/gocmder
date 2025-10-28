# ESC Key Navigation Test Guide

## Purpose
Verify that ESC key properly returns to Home page from all other pages, without exiting the application.

## Test Cases

### Test 1: ESC on Home Page
**Steps:**
1. Launch gocmder
2. Verify you're on Home page (F1)
3. Press ESC

**Expected Result:**
- Should remain on Home page
- Application should NOT exit
- No visible change

### Test 2: ESC from Terminal Page
**Steps:**
1. Press F2 to go to Terminal page
2. Verify terminal is displayed
3. Press ESC

**Expected Result:**
- Should return to Home page (F1)
- Help bar at bottom should update to remove "ESC Home" text

### Test 3: ESC from Database Page
**Steps:**
1. Press F4 to go to Database page
2. Verify database page is displayed
3. Press ESC

**Expected Result:**
- Should return to Home page (F1)
- Help bar should update

### Test 4: ESC from Tools Page
**Steps:**
1. Press F6 to go to Tools page
2. Verify tools list is displayed
3. Press ESC

**Expected Result:**
- Should return to Home page (F1)
- Help bar should update

### Test 5: ESC from Settings Page
**Steps:**
1. Press F7 to go to Settings page
2. Verify settings list is displayed
3. Press ESC

**Expected Result:**
- Should return to Home page (F1)
- Help bar should update

### Test 6: ESC from System Page
**Steps:**
1. Press F8 to go to System Information page
2. Verify system info is displayed
3. Press ESC

**Expected Result:**
- Should return to Home page (F1)
- Help bar should update

### Test 7: ESC with Dialog Open
**Steps:**
1. Press F4 to go to Database page
2. Press Ctrl+N to open connection dialog
3. Press ESC

**Expected Result:**
- Dialog should close
- Should remain on Database page
- Press ESC again to return to Home page

### Test 8: Multiple Navigation
**Steps:**
1. From Home, press F6 (Tools)
2. Press ESC (should go to Home)
3. Press F7 (Settings)
4. Press ESC (should go to Home)
5. Press F8 (System)
6. Press ESC (should go to Home)

**Expected Result:**
- Each ESC press returns to Home page
- No application exit at any point

## Verification Points
- ✓ ESC on Home page does not exit application
- ✓ ESC on any other page returns to Home page
- ✓ ESC on dialog closes dialog (not page)
- ✓ Help bar shows "ESC Home" on all pages except Home
- ✓ 'q' key still quits the application

## Implementation Details
- ESC key handling added in `globalInputHandler` in `internal/ui/app.go`
- Check for `SubDialogHasFocus()` ensures dialogs handle ESC first
- Only switches to Home if `currentPageIdx != homePageIndex`
- Help bar updated to show "ESC Home" on non-home pages

