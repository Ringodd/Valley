
from django.contrib import admin
from django.urls import path, include
from . import views
urlpatterns = [
   path('',views.index),
   path('main',views.userPage,name='main'),
   path('ajax/data/addPoints/',views.addPoints,name='add'),
   path('ajax/data/changeClass/',views.changeClass,name='add'),
   path('ajax/data/getList/',views.getList)
]
