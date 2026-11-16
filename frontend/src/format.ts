export const maskNumber = (n: string): string => (n.length <= 4 ? n : '••••' + n.slice(-4))

export const formatDate = (iso: string): string =>
  new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })