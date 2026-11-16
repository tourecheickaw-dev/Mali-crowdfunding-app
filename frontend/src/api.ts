const API: string = import.meta.env.VITE_API_URL
const DEV_ORGANIZER_ID: string | undefined = import.meta.env.VITE_DEV_ORGANIZER_ID

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const devHeaders: Record<string, string> = DEV_ORGANIZER_ID
    ? { 'X-Dev-Organizer-ID': DEV_ORGANIZER_ID }
    : {}
  const res = await fetch(API + path, {
    ...options,
    headers: { 'Content-Type': 'application/json', ...devHeaders, ...(options.headers ?? {}) },
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error((data as { error?: string }).error ?? 'An error occurred')
  }
  return data as T
}

export const fcfa = (n?: number): string =>
  new Intl.NumberFormat('fr-FR').format(n ?? 0) + ' FCFA'