import { readdir, readFile, stat } from "fs/promises";
import { join, resolve, relative } from "path";

export interface FileInfo {
  name: string;
  path: string;
  size: number;
}

export class LocalFileService {
  private readonly basePath: string;

  constructor(basePath: string = "./src/raccs") {
    this.basePath = resolve(basePath);
  }

  private getFullPath(filePath: string): string {
    const fullPath = resolve(join(this.basePath, filePath));
    
    // this fixes a path traversal vulnerability since we weren't validating this before
    const relativePath = relative(this.basePath, fullPath);
    if (relativePath.startsWith('..') || resolve(relativePath) !== relativePath) {
      throw new Error(`Access denied: Path traversal detected for ${filePath}`);
    }
    
    return fullPath;
  }

  async listFiles(
    subPath: string = "",
    extension?: string
  ): Promise<FileInfo[]> {
    try {
      const fullPath = this.getFullPath(subPath);
      const files = await readdir(fullPath);

      const fileInfos: FileInfo[] = [];

      for (const file of files) {
        const stats = await stat(join(fullPath, file));

        if (stats.isFile() && (!extension || file.endsWith(extension))) {
          fileInfos.push({
            name: file,
            path: join(subPath, file),
            size: stats.size,
          });
        }
      }

      return fileInfos.sort((a, b) => a.name.localeCompare(b.name));
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
