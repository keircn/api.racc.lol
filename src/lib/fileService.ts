import { readdir, readFile, stat } from "fs/promises";
import { join } from "path";

export interface FileInfo {
  name: string;
  path: string;
  size: number;
}

export class LocalFileService {
  constructor(private basePath: string = "./src/raccs") {}

  private getFullPath(filePath: string): string {
    return join(this.basePath, filePath);
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
