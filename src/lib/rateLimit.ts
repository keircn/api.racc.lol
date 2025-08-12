interface RateLimitEntry {
  count: number;
  resetTime: number;
}


class RateLimiter {
  private requests = new Map<string, RateLimitEntry>();
  private readonly maxRequests: number;
  private readonly windowMs: number;
  private lastCleanup: number = 0;
  private readonly cleanupInterval: number = 60000;

  constructor(maxRequests: number, windowMs: number) {
    this.maxRequests = maxRequests;
    this.windowMs = windowMs;
  }

  private getClientIP(request: Request): string {
    const cfConnectingIP = request.headers.get("CF-Connecting-IP");
    const xForwardedFor = request.headers.get("X-Forwarded-For");
    const xRealIP = request.headers.get("X-Real-IP");

    return (
      cfConnectingIP ||
      (xForwardedFor && xForwardedFor.split(",")[0].trim()) ||
      xRealIP ||
      "unknown"
    );
  }

  // run the cleanup every minute to be more proactive,
  // this avoids a memory leak when expired entries aren't
  // cleaned up and remain in memory indefinitely.
  // i also added a hard limit so it'll automatically
  // trim the fat if too many entries are logged

  private cleanupExpiredEntries(): void {
    const now = Date.now();
    
    if (now - this.lastCleanup < this.cleanupInterval) {
      return;
    }
    
    this.lastCleanup = now;
    
    for (const [key, entry] of this.requests.entries()) {
      if (now > entry.resetTime) {
        this.requests.delete(key);
      }
    }
    
    if (this.requests.size > 10000) {
      const sortedEntries = Array.from(this.requests.entries())
        .sort(([, a], [, b]) => a.resetTime - b.resetTime);
      
      const toKeep = sortedEntries.slice(-5000);
      this.requests.clear();
      
      for (const [key, entry] of toKeep) {
        this.requests.set(key, entry);
      }
    }
  }

  checkRateLimit(request: Request): {
    allowed: boolean;
    remaining: number;
    resetTime: number;
  } {
    const clientIP = this.getClientIP(request);
    const now = Date.now();

    this.cleanupExpiredEntries();

    const entry = this.requests.get(clientIP);

    if (!entry || now > entry.resetTime) {
      const resetTime = now + this.windowMs;
      this.requests.set(clientIP, { count: 1, resetTime });
      return {
        allowed: true,
        remaining: this.maxRequests - 1,
        resetTime,
      };
    }

    if (entry.count >= this.maxRequests) {
      return {
        allowed: false,
        remaining: 0,
        resetTime: entry.resetTime,
      };
    }

    entry.count++;
    return {
      allowed: true,
      remaining: this.maxRequests - entry.count,
      resetTime: entry.resetTime,
    };
  }
}

export const photoRateLimiter = new RateLimiter(5, 10 * 1000);
