import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { UpcomingEpisode } from '../../core/models/models';

@Component({
  selector: 'app-calendar',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './calendar.html',
  styleUrl: './calendar.scss'
})
export class Calendar implements OnInit {
  items = signal<UpcomingEpisode[]>([]);
  loading = signal(true);

  constructor(private media: MediaService) {}

  ngOnInit() {
    this.media.getCalendar().subscribe({
      next: (data) => {
        this.items.set(data);
        this.loading.set(false);
      },
      error: () => this.loading.set(false)
    });
  }
}