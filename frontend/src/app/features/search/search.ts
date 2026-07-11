import { Component, OnInit, OnDestroy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { MediaService } from '../../core/services/media.service';
import { SearchResult } from '../../core/models/models';
import { MediaModal } from '../../shared/media-modal/media-modal';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink, MediaModal],
  templateUrl: './search.html',
  styleUrl: './search.scss'
})
export class Search implements OnInit, OnDestroy {
  query = '';
  results = signal<SearchResult[]>([]);
  loading = signal(false);
  addedIds = signal<Set<string>>(new Set());
  skeletonItems = Array.from({ length: 8 });
  showingTrending = signal(false);
  selectedMediaId = signal<string | null>(null);

  private queryChanged = new Subject<string>();

  constructor(private media: MediaService, public auth: AuthService) {}

  ngOnInit() {
    this.refreshLibrary();
    this.loadTrending();

    this.queryChanged.pipe(
      debounceTime(400),
      distinctUntilChanged()
    ).subscribe((q) => {
      if (!q.trim()) {
        this.loadTrending();
      } else {
        this.doSearch(q);
      }
    });
  }

  ngOnDestroy() {
    this.queryChanged.complete();
  }

  refreshLibrary() {
    this.media.getProgress().subscribe({
      next: (items) => {
        const set = new Set(items.map(i => i.media_item_id));
        this.addedIds.set(set);
      }
    });
  }

  loadTrending() {
    this.loading.set(true);
    this.showingTrending.set(true);
    this.media.getTrending().subscribe({
      next: (data) => {
        this.results.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }

  doSearch(q: string) {
    this.loading.set(true);
    this.showingTrending.set(false);
    this.media.search(q).subscribe({
      next: (data) => {
        this.results.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }

  onQueryChange() {
    this.queryChanged.next(this.query);
  }

  submit() {
    if (!this.query.trim()) {
      this.loadTrending();
      return;
    }
    this.doSearch(this.query);
  }

  add(item: SearchResult) {
    this.media.addProgress(item.media_item_id, 'plan_to_watch').subscribe(() => {
      const set = new Set(this.addedIds());
      set.add(item.media_item_id);
      this.addedIds.set(set);
    });
  }

  isAdded(item: SearchResult): boolean {
    return this.addedIds().has(item.media_item_id);
  }

  openMedia(id: string) { this.selectedMediaId.set(id); }

  closeMedia() { this.selectedMediaId.set(null); }

  onMediaChanged() {
    this.refreshLibrary();
  }
}