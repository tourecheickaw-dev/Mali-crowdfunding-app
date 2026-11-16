import { useState } from 'react'
import { ImagePlus } from 'lucide-react'
import { uploadImage } from '../upload'

export interface PickedPhoto {
  url: string
  sensitive: boolean
}

export default function PhotoPicker({ onChange }: { onChange: (v: PickedPhoto | null) => void }) {
  const [url, setUrl] = useState('')
  const [sensitive, setSensitive] = useState(false)
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')

  async function onFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    if (!file) return
    setBusy(true)
    setError('')
    try {
      const publicUrl = await uploadImage(file)
      setUrl(publicUrl)
      onChange({ url: publicUrl, sensitive })
    } catch (err) {
      setUrl('')
      onChange(null)
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setBusy(false)
    }
  }

  function toggleSensitive(v: boolean) {
    setSensitive(v)
    if (url) onChange({ url, sensitive: v })
  }

  return (
    <div>
      <label className="flex h-48 cursor-pointer flex-col items-center justify-center gap-2 overflow-hidden rounded-3xl border-2 border-dashed border-gray-300 bg-white text-gray-500 hover:border-brand-500">
        {url ? (
          <img src={url} alt="" className={`h-full w-full object-cover ${sensitive ? 'blur-md' : ''}`} />
        ) : (
          <>
            <ImagePlus size={36} />
            <span className="text-lg font-medium">{busy ? 'Envoi en cours…' : 'Ajouter une photo'}</span>
          </>
        )}
        <input type="file" accept="image/jpeg,image/png,image/webp" className="hidden" onChange={onFile} />
      </label>

      {error && <p className="mt-2 text-red-600">{error}</p>}

      <label className="mt-3 flex items-start gap-3 text-brand-600rand-900ase">
        <input
          type="checkbox"
          checked={sensitive}
          onChange={(e) => toggleSensitive(e.target.checked)}
          className="mt-1 h-5 w-5"
        />
        <span>Cette photo montre un contenu médical ou difficile (elle sera floutée au départ)</span>
      </label>
      <p className="mt-2 text-sm text-gray-500">
        Merci de respecter la dignité des personnes qui apparaissent sur la photo.
      </p>
    </div>
  )
}