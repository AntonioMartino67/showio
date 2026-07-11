import { Component, OnInit, OnDestroy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { Subject } from 'rxjs';
import { debounceTime, distinctUntilChanged } from 'rxjs/operators';
import { MediaService } from '../../core/services/media.service';
import { SearchResult } from '../../core/models/models';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
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

  private queryChanged = new Subject<string>();

  constructor(private media: MediaService) {}

  ngOnInit() {
    // Carica la libreria dell'utente per sapere cosa è già in lista
    // (persiste tra ricariche, non solo nella sessione corrente)
    this.media.getProgress().subscribe({
      next: (items) => {
        const set = new Set(items.map(i => i.media_item_id));
        this.addedIds.set(set);
      }
    });

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
    // invio manuale (Invio o click sul bottone): salta il debounce
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
}