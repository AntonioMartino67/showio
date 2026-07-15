import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { AuthService } from '../../core/services/auth.service';
import { ProgressItem, Tag } from '../../core/models/models';
import { MediaModal } from '../../shared/media-modal/media-modal';
import { Loader } from '../../shared/loader/loader';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterLink, MediaModal, Loader],
  templateUrl: './dashboard.html',
  styleUrl: './dashboard.scss'
})
export class Dashboard implements OnInit {
  items = signal<ProgressItem[]>([]);
  loading = signal(true);
  filter = signal<'all' | 'watching' | 'completed' | 'dropped' | 'plan_to_watch'>('all');
  allTags = signal<Tag[]>([]);
  selectedTagId = signal<string | null>(null);

  constructor(private media: MediaService, public auth: AuthService) {}

  ngOnInit() {
    this.load();
    this.loadTags();
    this.auth.loadUser().subscribe();
  }

  loadTags() {
    this.media.getTags().subscribe({ next: (tags) => this.allTags.set(tags) });
  }

  load(silent = false) {
    if (!silent) this.loading.set(true);
    this.media.getProgress().subscribe({
      next: (data) => {
        this.items.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }

  filteredItems() {
    const f = this.filter();
    const tagId = this.selectedTagId();
    let result = f === 'all' ? this.items() : this.items().filter(i => i.status === f);
    if (tagId) {
      result = result.filter(i => i.tags?.some(t => t.id === tagId));
    }
    return result;
  }

  deleteTag(tagId: string, event: Event) {
    event.stopPropagation();
    if (!confirm('Eliminare questo tag? Verrà rimosso da tutti i titoli.')) return;
    this.media.deleteTag(tagId).subscribe(() => {
      if (this.selectedTagId() === tagId) this.selectedTagId.set(null);
      this.loadTags();
      this.load(true);
    });
  }

  onMediaChanged() {
    this.load(true);
    this.loadTags();
  }

  nextEpisode(item: ProgressItem) {
    const newEp = item.current_episode + 1;
    this.media.updateEpisode(item.media_item_id, item.current_season, newEp).subscribe(() => {
      this.load(true);
    });
  }

  selectedMediaId = signal<string | null>(null);
  openMedia(id: string) { this.selectedMediaId.set(id); }
  closeMedia() { this.selectedMediaId.set(null); }

  goToList() {
    this.filter.set('all');
    this.selectedTagId.set(null);
  }
}