import { api } from './api'

const ALLOWED = ['image/jpeg', 'image/png', 'image/webp']
export const MAX_IMAGE_BYTES = 5 * 1024 * 1024

// Reduce the photo before sending: faster on slow networks, and re-encoding removes GPS data.
export async function compressImage(file: File, maxSide = 1200, quality = 0.82): Promise<File> {
  if (!file.type.startsWith('image/')) return file
  const bmp = await createImageBitmap(file)
  const scale = Math.min(1, maxSide / Math.max(bmp.width, bmp.height))
  const w = Math.round(bmp.width * scale)
  const h = Math.round(bmp.height * scale)
  const canvas = document.createElement('canvas')
  canvas.width = w
  canvas.height = h
  const ctx = canvas.getContext('2d')
  if (!ctx) return file
  ctx.drawImage(bmp, 0, 0, w, h)
  bmp.close()
  const blob = await new Promise<Blob | null>((r) => canvas.toBlob(r, 'image/jpeg', quality))
  if (!blob) return file
  return new File([blob], file.name.replace(/\.\w+$/, '') + '.jpg', { type: 'image/jpeg' })
}

// Upload one image straight to S3 and return its public URL.
export async function uploadImage(original: File): Promise<string> {
  const file = await compressImage(original)
  if (!ALLOWED.includes(file.type)) throw new Error('Unsupported format: JPEG, PNG, or WebP')
  if (file.size > MAX_IMAGE_BYTES) throw new Error('Photo must be less than 5 MB')

  const p = await api<{ upload_url: string; public_url: string }>('/me/uploads/presign', {
    method: 'POST',
    body: JSON.stringify({ content_type: file.type, size: file.size }),
  })

  const res = await fetch(p.upload_url, {
    method: 'PUT',
    headers: { 'Content-Type': file.type },
    body: file,
  })
  if (!res.ok) throw new Error('Failed to upload photo')
  return p.public_url
}

// Attach an uploaded photo to a campaign as its cover.
export function attachCover(campaignId: string, url: string, sensitive: boolean) {
  return api(`/me/campaigns/${campaignId}/media`, {
    method: 'POST',
    body: JSON.stringify({ url, media_type: 'image', is_cover: true, is_sensitive: sensitive }),
  })
}