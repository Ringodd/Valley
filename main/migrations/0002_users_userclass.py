# Generated by Django 5.1.5 on 2025-01-24 09:34

import django.utils.timezone
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('main', '0001_initial'),
    ]

    operations = [
        migrations.AddField(
            model_name='users',
            name='userClass',
            field=models.TextField(default=django.utils.timezone.now, verbose_name='class'),
            preserve_default=False,
        ),
    ]
