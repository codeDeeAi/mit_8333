import type { ServiceRequest } from '@/types'
import { STATUS_LABELS, PRIORITY_LABELS } from './constants'

/** Trigger a client-side file download from a string blob. */
function download(filename: string, content: string, mime: string) {
  const blob = new Blob([content], { type: mime })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

function csvCell(value: string): string {
  const v = value ?? ''
  return /[",\n]/.test(v) ? `"${v.replace(/"/g, '""')}"` : v
}

/**
 * Export service requests to CSV. (The Go backend also exposes a
 * `/reports/requests.csv` endpoint; this is the client-side equivalent.)
 */
export function exportRequestsCsv(requests: ServiceRequest[]) {
  const headers = [
    'ID',
    'Title',
    'Category',
    'Location',
    'Priority',
    'Status',
    'Requested By',
    'Assigned Officer',
    'Created',
  ]
  const rows = requests.map((r) =>
    [
      r.id,
      r.title,
      r.categoryName,
      r.location,
      PRIORITY_LABELS[r.priority],
      STATUS_LABELS[r.status],
      r.createdByName,
      r.assignedOfficerName ?? '—',
      new Date(r.createdAt).toLocaleString(),
    ]
      .map((c) => csvCell(String(c)))
      .join(','),
  )
  download(
    `service-requests-${new Date().toISOString().slice(0, 10)}.csv`,
    [headers.join(','), ...rows].join('\n'),
    'text/csv;charset=utf-8',
  )
}

/**
 * "Export to PDF" via the browser print dialog (Save as PDF). Keeps the demo
 * dependency-free; the backend `/reports/requests.pdf` produces a server PDF.
 */
export function exportRequestsPdf(requests: ServiceRequest[]) {
  const win = window.open('', '_blank')
  if (!win) return
  const rows = requests
    .map(
      (r) => `<tr>
        <td>${r.title}</td>
        <td>${r.categoryName}</td>
        <td>${r.location}</td>
        <td>${PRIORITY_LABELS[r.priority]}</td>
        <td>${STATUS_LABELS[r.status]}</td>
        <td>${r.assignedOfficerName ?? '—'}</td>
      </tr>`,
    )
    .join('')
  win.document.write(`
    <html><head><title>Service Requests Report</title>
    <style>
      body { font-family: system-ui, sans-serif; padding: 32px; color: #111; }
      h1 { font-size: 20px; } p { color: #555; font-size: 12px; }
      table { width: 100%; border-collapse: collapse; margin-top: 16px; font-size: 12px; }
      th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
      th { background: #f4f4f4; }
    </style></head>
    <body>
      <h1>University Maintenance — Service Requests Report</h1>
      <p>Generated ${new Date().toLocaleString()} · ${requests.length} records</p>
      <table>
        <thead><tr><th>Title</th><th>Category</th><th>Location</th><th>Priority</th><th>Status</th><th>Officer</th></tr></thead>
        <tbody>${rows}</tbody>
      </table>
    </body></html>`)
  win.document.close()
  win.focus()
  setTimeout(() => win.print(), 300)
}
