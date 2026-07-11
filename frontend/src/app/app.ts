import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CommonModule } from '@angular/common';
import { inject as vercelInject } from '@vercel/analytics';
import { injectSpeedInsights } from '@vercel/speed-insights';
import { ErrorBannerService } from './core/services/error-banner.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, CommonModule],
  templateUrl: './app.html',
  styleUrl: './app.scss'
})
export class App {
  protected readonly title = signal('frontend');

  constructor(public banner: ErrorBannerService) {
    vercelInject();
    injectSpeedInsights();
  }
}