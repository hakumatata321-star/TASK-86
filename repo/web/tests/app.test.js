/**
 * app.test.js — unit tests for web/static/js/app.js
 *
 * app.js is a vanilla IIFE-based browser script that attaches helpers to
 * window (showToast, confirmAction) and sets up HTMX event listeners.
 * Tests run in jsdom via Vitest.
 */

import { describe, it, expect, beforeEach, vi } from 'vitest'
import { readFileSync } from 'fs'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const APP_JS = readFileSync(resolve(__dirname, '../static/js/app.js'), 'utf8')

// Re-run the script before each test so side-effects are reset.
function loadAppJs() {
  // Provide a minimal DOM structure the script expects.
  document.body.innerHTML = `
    <div id="toast-container"></div>
  `
  // eslint-disable-next-line no-eval
  eval(APP_JS)
}

// ---------------------------------------------------------------------------
// showToast
// ---------------------------------------------------------------------------

describe('showToast', () => {
  beforeEach(loadAppJs)

  it('appends a toast element to #toast-container', () => {
    window.showToast('Hello world', 'success')
    const container = document.getElementById('toast-container')
    expect(container.children.length).toBeGreaterThan(0)
  })

  it('includes the message text in the DOM', () => {
    window.showToast('Test message', 'info')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).toContain('Test message')
  })

  it('applies text-bg-success class for success type', () => {
    window.showToast('OK', 'success')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).toContain('text-bg-success')
  })

  it('applies text-bg-danger class for error type', () => {
    window.showToast('Error!', 'error')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).toContain('text-bg-danger')
  })

  it('applies text-bg-warning class for warning type', () => {
    window.showToast('Warn', 'warning')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).toContain('text-bg-warning')
  })

  it('applies text-bg-primary class for unknown type', () => {
    window.showToast('Note', 'unknown')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).toContain('text-bg-primary')
  })

  it('escapes HTML in the message to prevent XSS', () => {
    window.showToast('<script>alert("xss")</script>', 'error')
    const container = document.getElementById('toast-container')
    // Raw <script> tag must NOT appear in the rendered HTML.
    expect(container.innerHTML).not.toContain('<script>')
    expect(container.innerHTML).toContain('&lt;script&gt;')
  })

  it('escapes angle brackets and quotes', () => {
    window.showToast('<b>"quoted"</b>', 'info')
    const container = document.getElementById('toast-container')
    expect(container.innerHTML).not.toContain('<b>')
    expect(container.innerHTML).toContain('&lt;b&gt;')
  })

  it('does nothing when #toast-container is absent', () => {
    document.body.innerHTML = '' // remove container
    // Should not throw.
    expect(() => window.showToast('no container', 'info')).not.toThrow()
  })

  it('adds toast even when called multiple times', () => {
    window.showToast('First', 'success')
    window.showToast('Second', 'error')
    const container = document.getElementById('toast-container')
    expect(container.children.length).toBe(2)
  })
})

// ---------------------------------------------------------------------------
// confirmAction
// ---------------------------------------------------------------------------

describe('confirmAction', () => {
  beforeEach(loadAppJs)

  it('delegates to window.confirm', () => {
    const mockConfirm = vi.fn().mockReturnValue(true)
    window.confirm = mockConfirm

    const result = window.confirmAction('Delete this item?')

    expect(mockConfirm).toHaveBeenCalledWith('Delete this item?')
    expect(result).toBe(true)
  })

  it('returns false when user cancels', () => {
    window.confirm = vi.fn().mockReturnValue(false)
    expect(window.confirmAction('Sure?')).toBe(false)
  })

  it('uses default message when none provided', () => {
    const mockConfirm = vi.fn().mockReturnValue(true)
    window.confirm = mockConfirm

    window.confirmAction()

    expect(mockConfirm).toHaveBeenCalledWith('Are you sure?')
  })

  it('uses default message for empty string', () => {
    const mockConfirm = vi.fn().mockReturnValue(true)
    window.confirm = mockConfirm

    // The implementation: window.confirm(message || 'Are you sure?')
    window.confirmAction('')

    expect(mockConfirm).toHaveBeenCalledWith('Are you sure?')
  })
})

// ---------------------------------------------------------------------------
// Smoke test: script loads without throwing
// ---------------------------------------------------------------------------

describe('app.js smoke', () => {
  it('loads without throwing in jsdom environment', () => {
    document.body.innerHTML = '<div id="toast-container"></div>'
    expect(() => eval(APP_JS)).not.toThrow()
  })

  it('exports showToast on window after loading', () => {
    loadAppJs()
    expect(typeof window.showToast).toBe('function')
  })

  it('exports confirmAction on window after loading', () => {
    loadAppJs()
    expect(typeof window.confirmAction).toBe('function')
  })
})
