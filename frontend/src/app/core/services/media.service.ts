import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { SearchResult, ProgressItem, UpcomingEpisode, ProgressStatus, MediaDetail, Stats, Tag } from '../models/models';

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

  getMediaDetail(mediaItemId: string) {
    return this.http.get<MediaDetail>(`${environment.apiUrl}/media/${mediaItemId}`);
  }

  removeProgress(mediaItemId: string) {
    return this.http.delete(`${environment.apiUrl}/progress/${mediaItemId}`);
  }

  updateRating(mediaItemId: string, rating: number) {
    return this.http.put(`${environment.apiUrl}/progress/${mediaItemId}/rating`, { rating });
  }

  getTrending() {
    return this.http.get<SearchResult[]>(`${environment.apiUrl}/trending`);
  }

  getStats() {
    return this.http.get<Stats>(`${environment.apiUrl}/stats`);
  }

  getTags() {
    return this.http.get<Tag[]>(`${environment.apiUrl}/tags`);
  }

  createTag(name: string, color: string) {
    return this.http.post<Tag>(`${environment.apiUrl}/tags`, { name, color });
  }

  deleteTag(tagId: string) {
    return this.http.delete(`${environment.apiUrl}/tags/${tagId}`);
  }

  assignTag(mediaItemId: string, tagId: string) {
    return this.http.post(`${environment.apiUrl}/progress/${mediaItemId}/tags`, { tag_id: tagId });
  }

  removeTag(mediaItemId: string, tagId: string) {
    return this.http.delete(`${environment.apiUrl}/progress/${mediaItemId}/tags/${tagId}`);
  }
}