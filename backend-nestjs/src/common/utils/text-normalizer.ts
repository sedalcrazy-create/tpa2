/**
 * Normalizes Persian/Farsi text by converting Arabic characters to Persian equivalents
 * This ensures consistency in data storage and prevents issues with different character encodings
 */
export class TextNormalizer {
  /**
   * Converts Arabic characters (ي, ك) to Persian equivalents (ی, ک)
   * @param text - The text to normalize
   * @returns Normalized text with Persian characters
   */
  static normalize(text: string | null | undefined): string | null | undefined {
    if (!text) {
      return text;
    }

    return text
      .replace(/ي/g, 'ی') // Arabic Yeh to Persian Yeh
      .replace(/ك/g, 'ک'); // Arabic Kaf to Persian Kaf
  }

  /**
   * Normalizes multiple text fields in an object
   * @param obj - Object containing text fields to normalize
   * @param fields - Array of field names to normalize
   * @returns Object with normalized fields
   */
  static normalizeObject<T extends Record<string, any>>(
    obj: T,
    fields: (keyof T)[],
  ): T {
    const normalized = { ...obj };

    for (const field of fields) {
      if (typeof normalized[field] === 'string') {
        normalized[field] = this.normalize(normalized[field] as string) as T[keyof T];
      }
    }

    return normalized;
  }
}
