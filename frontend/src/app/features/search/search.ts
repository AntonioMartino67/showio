import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { SearchResult } from '../../core/models/models';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './search.html',
  styleUrl: './search.scss'
})
export class Search {
  query = '';
  results = signal<SearchResult[]>([]);
  loading = signal(false);
  addedIds = signal<Set<string>>(new Set());
  skeletonItems = Array.from({ length: 8 });

  constructor(private media: MediaService) {}

  submit() {
    if (!this.query.trim()) return;
    this.loading.set(true);
    this.media.search(this.query).subscribe({
      next: (data) => {
        this.results.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
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
}