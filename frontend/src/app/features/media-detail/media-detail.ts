import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { MediaDetail as MediaDetailData, ProgressStatus } from '../../core/models/models';

@Component({
  selector: 'app-media-detail',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './media-detail.html',
  styleUrl: './media-detail.scss'
})
export class MediaDetail implements OnInit {
  detail = signal<MediaDetailData | null>(null);
  loading = signal(true);
  mediaItemId = '';

  constructor(private route: ActivatedRoute, private router: Router, private media: MediaService) {}

  ngOnInit() {
    this.mediaItemId = this.route.snapshot.paramMap.get('id')!;
    this.load();
  }

  load() {
    this.loading.set(true);
    this.media.getMediaDetail(this.mediaItemId).subscribe({
      next: (data) => { this.detail.set(data); this.loading.set(false); },
      error: () => this.loading.set(false)
    });
  }

  addToList(status: ProgressStatus = 'plan_to_watch') {
    this.media.addProgress(this.mediaItemId, status).subscribe(() => this.load());
  }

  changeStatus(status: ProgressStatus) {
    this.media.addProgress(this.mediaItemId, status).subscribe(() => this.load());
  }

  markWatched(season: number, episode: number) {
    this.media.updateEpisode(this.mediaItemId, season, episode).subscribe(() => this.load());
  }

  setRating(rating: number) {
    this.media.updateRating(this.mediaItemId, rating).subscribe(() => this.load());
  }

  remove() {
    this.media.removeProgress(this.mediaItemId).subscribe(() => this.router.navigate(['/dashboard']));
  }

  isPast(d: MediaDetailData, season: number, episode: number): boolean {
    return season < d.current_season || (season === d.current_season && episode <= d.current_episode);
  }
}