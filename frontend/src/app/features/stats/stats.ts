import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { AuthService } from '../../core/services/auth.service';
import { Stats } from '../../core/models/models';
import { Loader } from '../../shared/loader/loader';
import { Navbar } from '../../shared/navbar/navbar';

@Component({
  selector: 'app-stats',
  standalone: true,
  imports: [CommonModule, RouterLink, Loader, Navbar],
  templateUrl: './stats.html',
  styleUrl: './stats.scss'
})
export class StatsPage implements OnInit {
  stats = signal<Stats | null>(null);
  loading = signal(true);

  constructor(private media: MediaService, public auth: AuthService) {}

  ngOnInit() {
    this.media.getStats().subscribe({
      next: (data) => { this.stats.set(data); this.loading.set(false); },
      error: () => this.loading.set(false)
    });
  }

  statusLabel(key: string): string {
    const map: Record<string, string> = {
      watching: 'In corso', completed: 'Completati', dropped: 'Abbandonati', plan_to_watch: 'Da vedere'
    };
    return map[key] || key;
  }

  maxStatusCount(): number {
    const s = this.stats();
    if (!s) return 1;
    return Math.max(1, ...Object.values(s.by_status));
  }

  statusEntries() {
    return Object.entries(this.stats()?.by_status || {});
  }

  typeEntries() {
    return Object.entries(this.stats()?.by_type || {});
  }
}