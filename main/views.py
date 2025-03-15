from django.shortcuts import render
from .models import Users
from .forms import LoginForm
from django.db import connection
from django.http import HttpResponseRedirect, HttpResponse
from django.http import JsonResponse
# Create your views here.
def index(request):
    users = Users.objects.all()
    error = None
    result = 'Incorrect'
    form = LoginForm()
    if request.method == 'POST':
        form = LoginForm(request.POST)
        
        if form.is_valid():
            username = request.POST['name']
            username = username.replace(' ','')
            password = request.POST['password']
            password = password.replace(' ','')
            print(username)
            print(password)
            usr = Users.authenicate(request,username=username,password=password)
            print(usr)
            if usr != None:
                print('suc')
                return HttpResponseRedirect('/main')
            else:
                print('err')
                error = 'Проверьте что ввели правильные данные!'
        else:
            print('err')
        
    context = {
        'data':users,
        'form':form,
        'result':result,
        'error':error
    }
    return render(request,'main/index.html',context)
def changeClass(request):
    if request.method == 'GET':
        print(request.GET.get('class'))

def addPoints(request):
    if request.method == 'GET':
        users = Users.objects.get(id = request.GET.get('userId'))
        users.count += int(request.GET.get('count'))
        if users.count < 0:
            users.count = 0
        #print(users.count)
        users.save()
        data = {'test':1}
        
        return JsonResponse(data)
    return JsonResponse({'error': 'Invalid request'}, status=400)
def getList(request):
    if request.method == 'GET':
        print(request)
        cursor = connection.cursor()
        students = request.GET.get('class')
        result = cursor.execute(f'SELECT * FROM main_users WHERE userClass = "{students}"').fetchall()
        print(result)
        data = {'students':result}
        
        return JsonResponse(data)
    return JsonResponse({'error': 'Invalid request'}, status=400)

def userPage(request):
    username = request.COOKIES.get('user')
    password = request.COOKIES.get('password')
    usr = Users.authenicate(request,username,password)
    cursor = connection.cursor()

    classes = cursor.execute('SELECT DISTINCT userClass from main_users WHERE userClass <> "null"').fetchall()
    print(classes)
    users = cursor.execute('SELECT * FROM main_users').fetchall()
    if usr != None:
        if usr.status == 1:
            data = {
            'name': usr.name,
            'password': usr.password,
            'count': usr.count,
            'status': usr.status
            }
            return render(request,'main/user.html',data)
        else:
            print(users)
            data = {
                'name': usr.name,
                'status': usr.status,
                'test': 'страница админа',
                'users':users,
                'classes':classes
            }
            return render(request, 'main/admin.html',data)
    
