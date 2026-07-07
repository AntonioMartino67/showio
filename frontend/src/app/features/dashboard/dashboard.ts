import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { AuthService } from '../../core/services/auth.service';
import { ProgressItem } from '../../core/models/models';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss'
})
export class Dashboard implements OnInit {
  items = signal<ProgressItem[]>([]);
  loading = signal(true);

  constructor(private media: MediaService, public auth: AuthService) {}

  ngOnInit() {
    this.load();
    this.auth.loadUser().subscribe();
  }

  load() {
    this.loading.set(true);
    this.media.getProgress().subscribe({
      next: (data) => {
        this.items.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }

  nextEpisode(item: ProgressItem) {
    const newEp = item.current_episode + 1;
    this.media.updateEpisode(item.media_item_id, item.current_season, newEp).subscribe(() => {
      this.load();
    });
  }
}