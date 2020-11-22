import { Component, OnInit } from '@angular/core';
import {BackendService} from '../../api/backend.service';

@Component({
  selector: 'app-home-page',
  templateUrl: './home-page.component.html',
  styleUrls: ['./home-page.component.css']
})
export class HomePageComponent implements OnInit {

  constructor(public backend: BackendService) { }

  ngOnInit(): void {
  }

}
