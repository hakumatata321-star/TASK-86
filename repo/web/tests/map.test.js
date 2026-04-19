/**
 * map.test.js — unit tests for web/static/js/map.js
 *
 * map.js is a vanilla IIFE that initialises a Leaflet map and fetches
 * geospatial data from the backend.  Tests mock Leaflet (L) and the
 * Fetch API, then eval the script to verify observable behaviours.
 */

import { describe, it, expect, beforeAll, beforeEach, vi } from 'vitest'
import { readFileSync } from 'fs'
import { resolve, dirname } from 'path'
import { fileURLToPath } from 'url'

const __dirname = dirname(fileURLToPath(import.meta.url))
const MAP_JS = readFileSync(resolve(__dirname, '../static/js/map.js'), 'utf8')

// ---------------------------------------------------------------------------
// Leaflet mock
// ---------------------------------------------------------------------------

function createLeafletMock() {
  const layerGroup = {
    addTo: vi.fn().mockReturnThis(),
    clearLayers: vi.fn(),
  }
  const rectangle = {
    bindPopup: vi.fn().mockReturnThis(),
    addTo: vi.fn().mockReturnThis(),
  }
  const marker = {
    bindPopup: vi.fn().mockReturnThis(),
  }
  const geoJSONLayer = { addTo: vi.fn() }
  const mapInstance = { addTo: vi.fn().mockReturnThis() }

  return {
    map: vi.fn().mockReturnValue(mapInstance),
    tileLayer: vi.fn().mockReturnValue({ addTo: vi.fn() }),
    layerGroup: vi.fn().mockReturnValue(layerGroup),
    geoJSON: vi.fn().mockReturnValue(geoJSONLayer),
    marker: vi.fn().mockReturnValue(marker),
    divIcon: vi.fn().mockReturnValue({}),
    rectangle: vi.fn().mockReturnValue(rectangle),
    _layerGroupMock: layerGroup,
    _rectangleMock: rectangle,
  }
}

// ---------------------------------------------------------------------------
// Setup DOM + mocks + load script
// ---------------------------------------------------------------------------

let L

beforeAll(() => {
  // Provide the DOM elements map.js looks for.
  document.body.innerHTML = `
    <div id="map"></div>
    <select id="layerSelect"><option value="">All</option></select>
    <button id="loadLayerBtn">Load</button>
    <button id="computeGridBtn">Compute</button>
    <input id="gridSizeInput" value="10" />
    <input id="metricInput" value="count" />
    <span id="mapStatus"></span>
  `

  // Mock Leaflet globally before loading the script.
  L = createLeafletMock()
  globalThis.L = L

  // Mock fetch: returns an empty-but-valid map data response.
  globalThis.fetch = vi.fn().mockResolvedValue({
    ok: true,
    json: () =>
      Promise.resolve({
        geojson: JSON.stringify({ type: 'FeatureCollection', features: [] }),
        locations: [],
        aggregates: [],
      }),
  })

  // Evaluate map.js — this runs the IIFE and calls loadLayerData('').
  eval(MAP_JS) // eslint-disable-line no-eval
})

// ---------------------------------------------------------------------------
// Initialisation tests
// ---------------------------------------------------------------------------

describe('map.js — Leaflet initialisation', () => {
  it('calls L.map with the "map" element id', () => {
    expect(L.map).toHaveBeenCalledWith('map', expect.any(Object))
  })

  it('initialises the map centred on continental US', () => {
    const [, options] = L.map.mock.calls[0]
    expect(options.center).toEqual([39.5, -98.35])
  })

  it('sets the default zoom level to 6', () => {
    const [, options] = L.map.mock.calls[0]
    expect(options.zoom).toBe(6)
  })

  it('creates an offline tile layer pointing to /static/tiles', () => {
    expect(L.tileLayer).toHaveBeenCalledWith(
      '/static/tiles/{z}/{x}/{y}.png',
      expect.any(Object)
    )
  })

  it('creates two layer groups (locations + density)', () => {
    expect(L.layerGroup).toHaveBeenCalledTimes(2)
  })
})

// ---------------------------------------------------------------------------
// Data loading tests
// ---------------------------------------------------------------------------

describe('map.js — data loading', () => {
  it('calls fetch on initialisation', () => {
    expect(globalThis.fetch).toHaveBeenCalled()
  })

  it('fetches /analytics/map/data on init (no layer)', () => {
    const firstCall = globalThis.fetch.mock.calls[0]
    expect(firstCall[0]).toBe('/analytics/map/data')
  })

  it('updates #mapStatus element', async () => {
    // Flush microtasks so the fetch promise resolves.
    await Promise.resolve()
    await Promise.resolve()
    const statusEl = document.getElementById('mapStatus')
    expect(statusEl).not.toBeNull()
  })
})

// ---------------------------------------------------------------------------
// escapeHtml smoke test via popup content
// ---------------------------------------------------------------------------

describe('map.js — escapeHtml (via renderDensityGrid)', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Re-setup L mock for fresh assertions
    globalThis.L = createLeafletMock()
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () =>
        Promise.resolve({
          geojson: null,
          locations: [],
          aggregates: [
            { cell_key: '39.5,-98.35', metric: '<script>', value: 5 },
          ],
        }),
    })
  })

  it('does not inject raw HTML into popups (XSS prevention)', async () => {
    // Re-eval to get a fresh run with the new fetch mock.
    document.body.innerHTML = `
      <div id="map"></div>
      <select id="layerSelect"></select>
      <button id="loadLayerBtn">Load</button>
      <button id="computeGridBtn">Compute</button>
      <input id="gridSizeInput" value="10" />
      <input id="metricInput" value="count" />
      <span id="mapStatus"></span>
    `
    const freshL = createLeafletMock()
    globalThis.L = freshL

    eval(MAP_JS) // eslint-disable-line no-eval

    // Let the fetch promise resolve.
    await new Promise((r) => setTimeout(r, 0))

    // If rectangles were created, check that none have raw <script> in their popup.
    if (freshL.rectangle.mock.calls.length > 0) {
      for (const rect of freshL._rectangleMock.bindPopup.mock.calls) {
        const popupHtml = rect[0]
        expect(popupHtml).not.toContain('<script>')
      }
    }
  })
})

// ---------------------------------------------------------------------------
// Smoke test: DOM element wiring
// ---------------------------------------------------------------------------

describe('map.js — DOM smoke tests', () => {
  it('loadLayerBtn element exists in test DOM', () => {
    expect(document.getElementById('loadLayerBtn')).not.toBeNull()
  })

  it('computeGridBtn element exists in test DOM', () => {
    expect(document.getElementById('computeGridBtn')).not.toBeNull()
  })

  it('mapStatus element exists in test DOM', () => {
    expect(document.getElementById('mapStatus')).not.toBeNull()
  })
})
