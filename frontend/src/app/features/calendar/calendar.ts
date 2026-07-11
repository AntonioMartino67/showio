import { Component, signal, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { MediaService } from '../../core/services/media.service';
import { AuthService } from '../../core/services/auth.service';
import { UpcomingEpisode } from '../../core/models/models';
import { Loader } from '../../shared/loader/loader';

@Component({
  selector: 'app-calendar',
  standalone: true,
  imports: [CommonModule, RouterLink, Loader],
  templateUrl: './calendar.html',
  styleUrl: './calendar.scss'
})
export class Calendar implements OnInit {
  items = signal<UpcomingEpisode[]>([]);
  loading = signal(true);

  constructor(private media: MediaService, public auth: AuthService) {}

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