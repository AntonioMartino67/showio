import { Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ErrorBannerService {
  message = signal<string | null>(null);
  private timer: any;

  show(msg: string) {
    this.message.set(msg);
    clearTimeout(this.timer);
    this.timer = setTimeout(() => this.message.set(null), 4000);
  }
}