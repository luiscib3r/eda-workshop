export function formatFileSize(bytes: number | string): string {
  // Convertir a número si es string
  const size = typeof bytes === "string" ? parseInt(bytes, 10) : bytes;

  if (isNaN(size) || size === 0) return "0 B";

  const units = ["B", "KB", "MB", "GB", "TB"];
  const k = 1024;
  const i = Math.floor(Math.log(size) / Math.log(k));

  return `${(size / Math.pow(k, i)).toFixed(2)} ${units[i]}`;
}
