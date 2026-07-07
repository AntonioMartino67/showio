import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { SearchResult, ProgressItem, UpcomingEpisode, ProgressStatus } from '../models/models';

@Injectable({ providedIn: 'root' })
export class MediaService {
  constructor(private http: HttpClient) {}

  search(query: string) {
    return this.http.get<SearchResult[]>(`${environment.apiUrl}/search`, {
      params: { q: query }
    });
  }

  getProgress() {
    return this.http.get<ProgressItem[]>(`${environment.apiUrl}/progress`);
  }

  addProgress(mediaItemId: string, status: ProgressStatus = 'plan_to_watch') {
    return this.http.post<{ id: string; status: string }>(`${environment.apiUrl}/progress`, {
      media_item_id: mediaItemId,
      status
    });
  }

  updateEpisode(mediaItemId: string, season: number, episode: number) {
    return this.http.put(`${environment.apiUrl}/progress/${mediaItemId}/episode`, {
      current_season: season,
      current_episode: episode
    });
  }

  getCalendar() {
    return this.http.get<UpcomingEpisode[]>(`${environment.apiUrl}/calendar`);
  }
}