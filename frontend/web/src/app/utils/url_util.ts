export class UrlUtil {
  public static appendSlash(str: string): string {
    return str.endsWith('/') ? str : str + '/'
  }
}
