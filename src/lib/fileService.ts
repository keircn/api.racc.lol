import { readdir, readFile, stat } from "fs/promises";
import { join, resolve, relative } from "path";

export interface FileInfo {
  name: string;
  path: string;
  size: number;
}

export class LocalFileService {
  private readonly basePath: string;

  constructor(basePath = "./src/raccs") {
    this.basePath = resolve(basePath);
  }

  private getFullPath(filePath: string): string {
    if (!filePath) {
      return this.basePath;
    }

    const fullPath = resolve(join(this.basePath, filePath));
    const relativePath = relative(this.basePath, fullPath);

    if (relativePath.startsWith("..") || relativePath.includes("..")) {
      throw new Error(`Access denied: Path traversal detected for ${filePath}`);
    }

    return fullPath;
  }

  async listFiles(subPath = "", extension?: string): Promise<FileInfo[]> {
    try {
      const fullPath = this.getFullPath(subPath);
      const files = await readdir(fullPath);

      const fileStats = await Promise.all(
        files.map(async (file) => {
          const filePath = join(fullPath, file);
          const stats = await stat(filePath);
          return { file, stats, isFile: stats.isFile() };
        })
      );

      return fileStats
        .filter(
          ({ isFile, file }) =>
            isFile && (!extension || file.endsWith(extension))
        )
        .map(({ file, stats }) => ({
          name: file,
          path: join(subPath, file),
          size: stats.size,
        }))
        .sort((a, b) => a.name.localeCompare(b.name));
    } catch (error) {
      console.error(`Error listing files in ${subPath}:`, error);
      return [];
    }
  }

  async getFile(filePath: string): Promise<Buffer | null> {
    try {
      return await readFile(this.getFullPath(filePath));
    } catch (error) {
      console.error(`Error reading file ${filePath}:`, error);
      return null;
    }
  }
}

export const fileService = new LocalFileService();
