import { Component, EventEmitter, Input, OnChanges, Output, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MediaService } from '../../core/services/media.service';
import { MediaDetail, ProgressStatus } from '../../core/models/models';
import { Loader } from '../loader/loader';

@Component({
  selector: 'app-media-modal',
  standalone: true,
  imports: [CommonModule, FormsModule, Loader],
  templateUrl: './media-modal.html',
  styleUrl: './media-modal.scss'
})
export class MediaModal implements OnChanges {
  @Input({ required: true }) mediaItemId!: string;
  @Output() closed = new EventEmitter<void>();
  @Output() changed = new EventEmitter<void>();

  detail = signal<MediaDetail | null>(null);
  loading = signal(true);

  constructor(private media: MediaService) {}

  ngOnChanges() { this.load(); }

  load() {
    this.loading.set(true);
    this.media.getMediaDetail(this.mediaItemId).subscribe({
      next: (data) => { this.detail.set(data); this.loading.set(false); },
      error: () => this.loading.set(false)
    });
  }

  close() { this.closed.emit(); }

  addToList(status: ProgressStatus = 'plan_to_watch') {
    this.media.addProgress(this.mediaItemId, status).subscribe(() => { this.load(); this.changed.emit(); });
  }

  changeStatus(status: ProgressStatus) {
    this.media.addProgress(this.mediaItemId, status).subscribe(() => { this.load(); this.changed.emit(); });
  }

  markWatched(season: number, episode: number) {
  const d = this.detail();
  if (!d) return;

  if (this.isPast(d, season, episode)) {
    // già segnato: lo smarchiamo tornando all'episodio precedente nella lista
    const idx = d.episodes.findIndex(e => e.season_number === season && e.episode_number === episode);
    const prev = idx > 0 ? d.episodes[idx - 1] : null;
    const prevSeason = prev ? prev.season_number : 0;
    const prevEpisode = prev ? prev.episode_number : 0;
    this.media.updateEpisode(this.mediaItemId, prevSeason, prevEpisode).subscribe(() => { this.load(); this.changed.emit(); });
  } else {
    this.media.updateEpisode(this.mediaItemId, season, episode).subscribe(() => { this.load(); this.changed.emit(); });
  }
}

  setRating(rating: number) {
    this.media.updateRating(this.mediaItemId, rating).subscribe(() => this.load());
  }

  remove() {
    this.media.removeProgress(this.mediaItemId).subscribe(() => { this.changed.emit(); this.close(); });
  }

  isPast(d: MediaDetail, season: number, episode: number): boolean {
    return season < d.current_season || (season === d.current_season && episode <= d.current_episode);
  }
}